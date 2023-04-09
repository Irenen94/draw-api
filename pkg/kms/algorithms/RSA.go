package algorithms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type RSA struct {
}

func NewRsaCipher() *RSA {
	return &RSA{}
}

func (c *RSA) Encrypt(privateKey, plaintext []byte) ([]byte, error) {

	rsaPrivateKey, err := BytesToPrivateKeyRsa(privateKey)
	if err != nil {
		return nil, err
	}
	rsaPubKey := &rsaPrivateKey.PublicKey

	label := []byte("OAEP Encrypted")
	k := rsaPubKey.Size()
	DefaultMaxEncryptBlock := k - 2*sha256.New().Size() - 2

	var msg []byte
	dst := make([]byte, 0)
	for {

		if len(plaintext) > DefaultMaxEncryptBlock {
			msg = plaintext[:DefaultMaxEncryptBlock]
			plaintext = plaintext[DefaultMaxEncryptBlock:]
			ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, msg, label)
			if err != nil {
				return nil, err
			}
			dst = append(dst, ciphertext...)
		} else {
			msg = plaintext[:]
			ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, msg, label)
			if err != nil {
				return nil, err
			}
			dst = append(dst, ciphertext...)
			break
		}

	}
	return dst, nil
}
func (c *RSA) Decrypt(privateKey, in []byte) ([]byte, error) {

	rsaPrivateKey, err := BytesToPrivateKeyRsa(privateKey)
	if err != nil {
		return nil, err
	}

	label := []byte("OAEP Encrypted")

	var cipherBlob []byte
	DefaultMaxDecryptBlock := rsaPrivateKey.Size()
	dst := make([]byte, 0)
	for {
		if len(in) > DefaultMaxDecryptBlock {
			cipherBlob = in[:DefaultMaxDecryptBlock]
			in = in[DefaultMaxDecryptBlock:]
			out, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, cipherBlob, label)
			if err != nil {
				return nil, err
			}
			dst = append(dst, out...)
		} else {
			cipherBlob = in[:]
			out, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, cipherBlob, label)
			if err != nil {
				return nil, err
			}
			dst = append(dst, out...)
			break
		}
	}
	return dst, nil
}
