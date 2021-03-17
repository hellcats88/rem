package pkcs11

import (
	"github.com/ThalesIgnite/crypto11"
	"github.com/hellcats88/abstracte/security"
)

// NewCloudHSM creates an instance of HSM Secure module used to sign and crypt information based on AWS CloudHSM
func NewCloudHSM(config HSMConfig) (security.SecureModule, error) {
	hsmCtx, err := crypto11.Configure(&crypto11.Config{
		Path:            config.Path,
		Pin:             config.Pin,
		TokenLabel:      config.Label,
		UseGCMIVFromHSM: true,
	})

	if err != nil {
		return nil, err
	}

	return hsm{
		cryptoCtx: hsmCtx,
	}, nil
}
