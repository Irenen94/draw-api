package kms

import (
	"aiot-service-for-mfp/pkg/kms/algorithms/aead"
	"aiot-service-for-mfp/pkg/kms/cache"
	"aiot-service-for-mfp/pkg/kms/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type RemoteKmsClient struct {
	client *Client
}

func NewRemoteKmsClient(config *config.LongforConfig) *RemoteKmsClient {
	return &RemoteKmsClient{
		client: NewClientWithAccessKey(config),
	}
}

func (k *RemoteKmsClient) Encrypt(request *EncryptRequest) (*EncryptResponse, error) {
	return k.client.Encrypt(request)
}

func (k *RemoteKmsClient) Decrypt(request *DecryptRequest) (*DecryptResponse, error) {
	return k.client.Decrypt(request)
}
func (k *RemoteKmsClient) NewSecret(request *NewSecretRequest) (*NewSecretResponse, error) {
	return k.client.NewSecret(request)
}
func (k *RemoteKmsClient) SetSecret(request *SetSecretValueRequest) (*SetSecretValueResponse, error) {
	return k.client.SetSecretValue(request)
}
func (k *RemoteKmsClient) GetSecret(request *GetSecretValueRequest) (*GetSecretValueResponse, error) {
	return k.client.GetSecret(request)
}
func (k *RemoteKmsClient) getKeyInfo(keyid, keyVersionId string) (*KeyInfo, error) {
	var (
		err     error
		keyInfo = KeyInfo{}
	)

	cacheKey := fmt.Sprintf("%s%s", keyid, keyVersionId)
	if err = cache.Get(cacheKey, &keyInfo); err == nil {
		return &keyInfo, nil
	}

	var rander = rand.Reader
	privateKeyRSA, err := rsa.GenerateKey(rander, 2048)
	if err != nil {
		return nil, err
	}

	publicKey := privateKeyRSA.PublicKey
	pkiPublicKey := x509.MarshalPKCS1PublicKey(&publicKey)

	publicKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkiPublicKey,
	})

	request := CreateExportKeySecretRequest()
	request.SetPublicKey(string(publicKeyBytes))
	request.SetKeyId(keyid)
	request.SetKeyVersionId(keyVersionId)

	resp, err := k.client.ExportKeySecret(request)
	if err != nil {
		return nil, err
	}

	cipherBlob, err := base64.StdEncoding.DecodeString(resp.SecretKey)
	if err != nil {
		return nil, err
	}

	if len(cipherBlob) > 32 {
		cipherBlob = cipherBlob[32:]
	}

	label := []byte("OAEP Encrypted")
	secretKey, err := rsa.DecryptOAEP(sha256.New(), rander, privateKeyRSA, cipherBlob, label)
	if err != nil {
		return nil, err
	}

	aw := aead.Wrapper{}
	err = aw.SetAESKeyBytes(secretKey)
	if err != nil {
		return nil, err
	}

	in, err := base64.StdEncoding.DecodeString(resp.Secret)
	if err != nil {
		return nil, err
	}
	if len(in) > 32 {
		in = in[32:]
	}
	out, err := aw.Decrypt(in, nil)
	if err != nil {
		return nil, err
	}

	val := NewKeyInfo(resp.KeyId, resp.KeySpec, resp.KeyUsage, resp.State, resp.KeyVersionId, out)
	_ = cache.Set(cacheKey, val, k.client.config.GetLocalCacheTimeout())
	return &val, nil
}
