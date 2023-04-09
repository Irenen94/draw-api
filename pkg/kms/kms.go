package kms

import (
	"aiot-service-for-mfp/pkg/kms/config"
)

type IKmsClient interface {
	Encrypt(request *EncryptRequest) (*EncryptResponse, error)
	Decrypt(request *DecryptRequest) (*DecryptResponse, error)
	NewSecret(request *NewSecretRequest) (*NewSecretResponse, error)
	GetSecret(request *GetSecretValueRequest) (*GetSecretValueResponse, error)
	SetSecret(request *SetSecretValueRequest) (*SetSecretValueResponse, error)
}

type KmsClient struct {
	config config.LongforConfig
}

func NewKmsClient(config *config.LongforConfig) IKmsClient {
	if config.GetKmsMode().String() == "LOCAL" {
		return NewLocalKmsClient(config)
	}
	return NewRemoteKmsClient(config)
}
