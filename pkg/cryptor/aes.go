package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

type AesCryptor struct {
	Key []byte
	Iv  []byte
}

func NewAesCryptor(base64key string, base64iv string) AesCryptor {
	key, _ := base64.StdEncoding.DecodeString(base64key)
	iv, _ := base64.StdEncoding.DecodeString(base64iv)

	return AesCryptor{
		Key: key,
		Iv:  iv,
	}
}

func (a *AesCryptor) AESEncrypt(src string) []byte {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(a.Iv))
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return crypted
}

func (a *AesCryptor) AESDecrypt(crypt []byte) (result []byte, err error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		fmt.Println("key error1", err)
		return []byte{}, err
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
		return []byte{}, nil
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(a.Iv))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)

	return PKCS5Trimming(decrypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
