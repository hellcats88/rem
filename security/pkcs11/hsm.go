package pkcs11

import (
	"crypto"
	"crypto/cipher"
	"fmt"

	"github.com/ThalesIgnite/crypto11"
	"github.com/hellcats88/abstracte/runtime"
	"github.com/hellcats88/abstracte/security"
)

type hsm struct {
	cryptoCtx              *crypto11.Context
	symmetricCapabilities  map[string][]int
	asymmetricCapabilities map[string][]int
}

// HSMConfig setup the HSM library to be used
type HSMConfig struct {
	Path                   string
	Pin                    string
	Label                  string
	SymmetricCapabilities  map[string][]int
	AsymmetricCapabilities map[string][]int
}

// New creates an instance of HSM Secure module used to sign and crypt information
func New(config HSMConfig) (security.SecureModule, error) {
	hsmCtx, err := crypto11.Configure(&crypto11.Config{
		Path:       config.Path,
		Pin:        config.Pin,
		TokenLabel: config.Label,
	})

	if err != nil {
		return nil, err
	}

	if config.SymmetricCapabilities == nil {
		config.SymmetricCapabilities["AES"] = []int{128}
	}
	if config.AsymmetricCapabilities == nil {
		config.SymmetricCapabilities["RSA"] = []int{1024}
	}

	return hsm{
		cryptoCtx:              hsmCtx,
		symmetricCapabilities:  config.SymmetricCapabilities,
		asymmetricCapabilities: config.AsymmetricCapabilities,
	}, nil
}

func (mod hsm) GenerateRSAKeyPair(ctx runtime.Context, req security.GenerateRSAKeyPairReq) (crypto.PublicKey, error) {
	key, err := mod.cryptoCtx.GenerateRSAKeyPairWithLabel([]byte(req.Alias), []byte(req.Alias), req.Bits)
	if err != nil {
		return nil, err
	}

	return key.Public(), nil
}

func (mod hsm) GenerateECDSAKeyPair(ctx runtime.Context, req security.GenerateECDSAKeyPairReq) (crypto.PublicKey, error) {
	key, err := mod.cryptoCtx.GenerateECDSAKeyPairWithLabel([]byte(req.Alias), []byte(req.Alias), req.Bits)
	if err != nil {
		return nil, err
	}

	return key.Public(), nil
}

func (mod hsm) GenerateAESKey(ctx runtime.Context, req security.GenerateAESKeyReq) (cipher.Block, error) {
	return mod.cryptoCtx.GenerateSecretKeyWithLabel([]byte(req.Alias), []byte(req.Alias), req.Bits, crypto11.CipherAES)
}

func (mod hsm) Signer(ctx runtime.Context, alias string) (crypto.Signer, error) {
	key, err := mod.cryptoCtx.FindKeyPair([]byte(alias), []byte(alias))
	if err != nil {
		return nil, err
	}

	if key == nil {
		return nil, fmt.Errorf("Search for key %s in HSM failed. Key not found", alias)
	}

	return key, nil
}

func (mod hsm) Block(ctx runtime.Context, alias string) (cipher.Block, error) {
	block, err := mod.cryptoCtx.FindKey(nil, []byte(alias))
	if err != nil {
		return nil, err
	}

	if block == nil {
		return nil, fmt.Errorf("Search for key %s in HSM failed. Key not found", alias)
	}

	return block, nil
}

func (mod hsm) Capabilities() security.CapabilitiesResp {
	var symmetricCaps []security.CapabilityResp
	for capName, capLen := range mod.symmetricCapabilities {
		symmetricCaps = append(symmetricCaps, security.CapabilityResp{
			Name: capName,
			Len:  capLen,
		})
	}

	var asymmetricCaps []security.CapabilityResp
	for capName, capLen := range mod.asymmetricCapabilities {
		asymmetricCaps = append(asymmetricCaps, security.CapabilityResp{
			Name: capName,
			Len:  capLen,
		})
	}

	return security.CapabilitiesResp{
		Asymmetric: asymmetricCaps,
		Symmetric:  symmetricCaps,
	}
}
