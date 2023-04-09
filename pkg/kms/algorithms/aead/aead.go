package aead

import (
	"aiot-service-for-mfp/pkg/kms/uuid"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// Wrapper implements the wrapping.Wrapper interface for AEAD
type Wrapper struct {
	keyID    string
	keyBytes []byte
	aead     cipher.AEAD
}
type EncryptedBlobInfo struct {
	Ciphertext []byte
	keyID      string
}

//关键参数应该是AES密钥，16，24或32个字节来选择 AES-128，AES-192 或 AES-256。

func GenerateAESKey(size int) (string, error) {
	if size == 16 || size == 24 || size == 32 {
		keyRaw, err := uuid.GenerateRandomBytes(size)
		if err != nil {
			return "", err
		}
		key := base64.StdEncoding.EncodeToString(keyRaw)
		return key, nil
	}
	return "", errors.New("invalid aes key size")
}
func GenerateAESKeyBytes(size int) ([]byte, error) {
	if size == 16 || size == 24 || size == 32 {
		return uuid.GenerateRandomBytes(size)
	}
	return nil, errors.New("invalid aes key size")
}
func (s *Wrapper) SetAESKeyString(key string) error {
	keyRaw, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	return s.SetAESKeyBytes(keyRaw)
}
func (s *Wrapper) SetAESKeyBytes(key []byte) error {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aead, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return err
	}

	s.keyBytes = key
	s.aead = aead
	return nil
}
func (s *Wrapper) Encrypt(plaintext, aad []byte) ([]byte, error) {

	iv, err := uuid.GenerateRandomBytes(12)
	if err != nil {
		return nil, err
	}
	ciphertext := s.aead.Seal(nil, iv, plaintext, aad)
	Ciphertext := append(iv, ciphertext...)
	return Ciphertext, nil
}
func (s *Wrapper) EncryptWithZeroIV(plaintext, aad []byte) ([]byte, error) {
	iv := make([]byte, 12)
	ciphertext := s.aead.Seal(nil, iv, plaintext, aad)
	Ciphertext := append(iv, ciphertext...)
	return Ciphertext, nil
}
func (s *Wrapper) Decrypt(in, aad []byte) ([]byte, error) {

	iv, ciphertext := in[:12], in[12:]
	plaintext, err := s.aead.Open(nil, iv, ciphertext, aad)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}
