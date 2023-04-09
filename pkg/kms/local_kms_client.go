package kms

import (
	"aiot-service-for-mfp/pkg/kms/algorithms"
	"aiot-service-for-mfp/pkg/kms/config"
	"aiot-service-for-mfp/pkg/kms/uuid"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type LocalKmsClient struct {
	*RemoteKmsClient
}

func NewLocalKmsClient(config *config.LongforConfig) *LocalKmsClient {
	return &LocalKmsClient{
		RemoteKmsClient: NewRemoteKmsClient(config),
	}
}
func (k *LocalKmsClient) Encrypt(request *EncryptRequest) (*EncryptResponse, error) {
	key, err := k.getKeyInfo(request.GetKeyId(), "")
	if err != nil {
		return nil, err
	}

	var out []byte

	plaintext, err := base64.StdEncoding.DecodeString(request.GetPlaintext())
	if err != nil {
		return nil, err
	}
	if strings.Contains(key.KeySpec, "AES") {
		out, err = algorithms.NewAesCipher().Encrypt(key.Secret, plaintext, request.GetRandom())
	} else if strings.Contains(key.KeySpec, "RSA") {
		out, err = algorithms.NewRsaCipher().Encrypt(key.Secret, plaintext)
	} else if strings.Contains(key.KeySpec, "LIP002") {
		out, err = algorithms.NewTripleEcbDes().Encrypt(key.Secret, plaintext)
	} else {
		return nil, errors.New(fmt.Sprintf("KeySpec not support.%s", key.KeySpec))
	}
	keyIdByte, err := uuid.ParseUUID(key.KeyId)
	keyVersionByte, err := uuid.ParseUUID(key.KeyVersionId)

	encrypted := append(keyIdByte, keyVersionByte...)
	encrypted = append(encrypted, out...)

	requestId, _ := uuid.GenerateUUID()

	return &EncryptResponse{
		KeyId:          key.KeyId,
		KeyVersionId:   key.KeyVersionId,
		CiphertextBlob: base64.StdEncoding.EncodeToString(encrypted),
		RequestId:      requestId,
	}, nil
}
func (k *LocalKmsClient) Decrypt(request *DecryptRequest) (*DecryptResponse, error) {

	ciphertextBlob, err := base64.StdEncoding.DecodeString(request.ciphertextBlob)
	if err != nil {
		return nil, err
	}

	keyId, err := uuid.FormatUUID(ciphertextBlob[0:16])
	if err != nil {
		return nil, err
	}
	keyVer, err := uuid.FormatUUID(ciphertextBlob[16:32])
	if err != nil {
		return nil, err
	}

	key, err := k.getKeyInfo(keyId, keyVer)
	if err != nil {
		return nil, err
	}

	var in, out []byte
	in = ciphertextBlob[32:]

	if strings.Contains(key.KeySpec, "AES") {
		out, err = algorithms.NewAesCipher().Decrypt(key.Secret, in)
	} else if strings.Contains(key.KeySpec, "RSA") {
		out, err = algorithms.NewRsaCipher().Decrypt(key.Secret, in)
	} else if strings.Contains(key.KeySpec, "LIP002") {
		out, err = algorithms.NewTripleEcbDes().Decrypt(key.Secret, in)
	} else {
		return nil, errors.New(fmt.Sprintf("KeySpec not support.%s", key.KeySpec))
	}

	requestId, _ := uuid.GenerateUUID()

	return &DecryptResponse{
		KeyId:        key.KeyId,
		KeyVersionId: key.KeyVersionId,
		Plaintext:    base64.StdEncoding.EncodeToString(out),
		RequestId:    requestId,
	}, nil
}
