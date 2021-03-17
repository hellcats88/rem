package kms

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/x509"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/hellcats88/abstracte/runtime"
	"github.com/hellcats88/abstracte/security"
	"github.com/lstoll/awskms"
)

// KMSStorageData is the entity stored in the KMS Storage
type KMSStorageData struct {
	Key  []byte
	Algo string
}

// KMSStorage defines the behavior of the KMS storage.
// KMS needs a master KEY (stored in the AWS KMS service) used to generate a random AES key
// This interface will be used to store the random AES Key that will be used to encrypt or decrypt data.
// This AES Key is encrypted with the master Key and can be used only with Encrypt/Decrypt KMS client operations
type KMSStorage interface {
	Save(ctx runtime.Context, alias string, criptedKey []byte, algo string) error
	Get(ctx runtime.Context, alias string) (KMSStorageData, error)
}

type kmsmod struct {
	kmsClient               kmsiface.KMSAPI
	storage                 KMSStorage
	symmetricMasterKeyAlias string
}

var rsaTypes = map[int]string{
	2048: kms.CustomerMasterKeySpecRsa2048,
	3072: kms.CustomerMasterKeySpecRsa3072,
	4096: kms.CustomerMasterKeySpecRsa4096,
}

var ecdsaTypes = map[elliptic.Curve]string{
	elliptic.P256(): kms.CustomerMasterKeySpecEccNistP256,
	elliptic.P384(): kms.CustomerMasterKeySpecEccNistP384,
	elliptic.P521(): kms.CustomerMasterKeySpecEccNistP521,
}

// New creates and instance of secure module based on AWS KMS service
// only for asymmetric security
func NewAsync(kmsClient kmsiface.KMSAPI) security.SecureModule {
	return kmsmod{
		kmsClient: kmsClient,
	}
}

// New creates and instance of secure module based on AWS KMS service
// for asymmetric and symmetric security
func New(kmsClient kmsiface.KMSAPI, storage KMSStorage, symmetricMasterKeyAlias string) security.SecureModule {
	return kmsmod{
		kmsClient:               kmsClient,
		storage:                 storage,
		symmetricMasterKeyAlias: symmetricMasterKeyAlias,
	}
}

func (mod kmsmod) generateKey(keyType string, keyAlias string) (crypto.PublicKey, error) {
	usage := kms.KeyUsageTypeEncryptDecrypt
	key, err := mod.kmsClient.CreateKey(&kms.CreateKeyInput{
		CustomerMasterKeySpec: &keyType,
		KeyUsage:              &usage,
		Description:           &keyAlias,
	})

	if err != nil {
		return nil, err
	}

	_, aliasErr := mod.kmsClient.CreateAlias(&kms.CreateAliasInput{
		AliasName:   &keyAlias,
		TargetKeyId: key.KeyMetadata.KeyId,
	})

	if aliasErr != nil {
		return nil, aliasErr
	}

	pubKey, pubKeyErr := mod.kmsClient.GetPublicKey(&kms.GetPublicKeyInput{KeyId: &keyAlias})
	if pubKeyErr != nil {
		return nil, pubKeyErr
	}

	return x509.ParsePKIXPublicKey(pubKey.PublicKey)
}

func (mod kmsmod) generateSymmetricKey(ctx runtime.Context, alias string, len int) (cipher.Block, error) {
	defaultReq := kms.GenerateDataKeyInput{
		KeyId: &mod.symmetricMasterKeyAlias,
	}

	if len == 128 {
		algo := kms.DataKeySpecAes128
		defaultReq.KeySpec = &algo
	} else if len == 256 {
		algo := kms.DataKeySpecAes256
		defaultReq.KeySpec = &algo
	} else {
		customLen := int64(len)
		defaultReq.NumberOfBytes = &customLen
	}

	dataKeyResp, err := mod.kmsClient.GenerateDataKey(&defaultReq)
	if err != nil {
		return nil, err
	}

	if err := mod.storage.Save(ctx, alias, dataKeyResp.CiphertextBlob, *defaultReq.KeySpec); err != nil {
		return nil, err
	}

	return aes.NewCipher(dataKeyResp.Plaintext)
}

func (mod kmsmod) GenerateRSAKeyPair(ctx runtime.Context, req security.GenerateRSAKeyPairReq) (crypto.PublicKey, error) {
	return mod.generateKey(rsaTypes[req.Bits], req.Alias)
}

func (mod kmsmod) GenerateECDSAKeyPair(ctx runtime.Context, req security.GenerateECDSAKeyPairReq) (crypto.PublicKey, error) {
	return mod.generateKey(ecdsaTypes[req.Bits], req.Alias)
}

func (mod kmsmod) GenerateAESKey(ctx runtime.Context, req security.GenerateAESKeyReq) (cipher.Block, error) {
	return mod.generateSymmetricKey(ctx, req.Alias, req.Bits)
}

func (mod kmsmod) Signer(ctx runtime.Context, alias string) (crypto.Signer, error) {
	return awskms.NewSigner(context.Background(), mod.kmsClient, alias)
}

func (mod kmsmod) Block(ctx runtime.Context, alias string) (cipher.Block, error) {
	key, err := mod.storage.Get(ctx, alias)
	if err != nil {
		return nil, err
	}

	descrResp, err := mod.kmsClient.Decrypt(&kms.DecryptInput{KeyId: &mod.symmetricMasterKeyAlias, EncryptionAlgorithm: &key.Algo, CiphertextBlob: key.Key})
	if err != nil {
		return nil, err
	}

	return aes.NewCipher(descrResp.Plaintext)
}

func (mod kmsmod) Capabilities() security.CapabilitiesResp {
	return security.CapabilitiesResp{
		Asymmetric: []security.CapabilityResp{
			{Name: "ECDSA", Len: []int{256, 384, 521}},
		},
		Symmetric: []security.CapabilityResp{
			{Name: "AES", Len: []int{128, 192, 256}},
		},
	}
}
