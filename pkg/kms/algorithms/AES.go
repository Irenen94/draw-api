package algorithms

import "aiot-service-for-mfp/pkg/kms/algorithms/aead"

type AES struct {
}

func NewAesCipher() *AES {
	return &AES{}
}
func (c *AES) Encrypt(privateKey, plaintext []byte, random bool) ([]byte, error) {
	aw := aead.Wrapper{}
	if err := aw.SetAESKeyBytes([]byte(privateKey)); err != nil {
		return nil, err
	}
	if random {
		return aw.Encrypt([]byte(plaintext), nil)
	}
	return aw.EncryptWithZeroIV([]byte(plaintext), nil)
}
func (c *AES) Decrypt(privateKey, in []byte) ([]byte, error) {
	aw := aead.Wrapper{}
	err := aw.SetAESKeyBytes(privateKey)
	if err != nil {
		return nil, err
	}
	out, err := aw.Decrypt(in, nil)
	if err != nil {
		return nil, err
	}
	return out, nil
}
