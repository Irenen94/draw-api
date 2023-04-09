package algorithms

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

var RandReader = rand.Reader

func GenerateKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := io.ReadFull(RandReader, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
func GenerateAESKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := io.ReadFull(RandReader, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}
func GeneratePrivateKeyRsa(keyBits int) ([]byte, error) {
	privateKey, err := rsa.GenerateKey(RandReader, keyBits)
	if err != nil {
		return nil, err
	}
	privateKeyRsa := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return privateKeyRsa, nil
}

func GenerateECKey(keyBits int) ([]byte, error) {
	var curve elliptic.Curve
	switch keyBits {
	case 224:
		curve = elliptic.P224()
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, errors.New(fmt.Sprintf("unsupported EC key:%d", keyBits))

	}
	privateKey, err := ecdsa.GenerateKey(curve, RandReader)
	if err != nil {
		return nil, err
	}
	//MarshalECPrivateKey marshals an EC private key into ASN.1, DER format.
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	if privateKeyBytes == nil {
		return nil, errors.New("no data returned when marshalling to private key")
	}

	privateKeyEcc := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return privateKeyEcc, nil
}

// GeneratePrivateKey /*  support AES_128,AES_256,RAS_2048,RAS_3072,RSA_4096,EC_P224,EC_P256,EC_P384,EC_P521
func GeneratePrivateKey(keyType string, keyBits int) ([]byte, error) {
	switch keyType {
	case "AES":
		return GenerateAESKey(keyBits / 8)
	case "RSA":
		return GeneratePrivateKeyRsa(keyBits)
	case "EC":
		return GenerateECKey(keyBits)
	case "FPE":
		//由于fpe底层加密和AES一致，因此直接使用AES密钥生成算法
		return GenerateAESKey(keyBits / 8)
	default:
		return nil, fmt.Errorf("not support %s_%d", keyType, keyBits)
	}
}

func BytesToPrivateKeyRsa(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	b := block.Bytes
	if block.Type == "RSA PRIVATE KEY" {
		return x509.ParsePKCS1PrivateKey(b)
	} else {
		key, err := x509.ParsePKCS8PrivateKey(b)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("failed parse RSA PrivateKey")
		}
		return rsaKey, nil
	}
}
func RsaKeyToByte(privateKey *rsa.PrivateKey) ([]byte, error) {
	keyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return keyBytes, nil
}
func EcKeyToByte(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	ecBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	keyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: ecBytes,
	})
	return keyBytes, nil
}
func BytesToPrivateKeyEcc(priv []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	b := block.Bytes
	var err error
	key, err := x509.ParseECPrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}
func BytesToPublicKeyRsa(pbk []byte) (*rsa.PublicKey, error) {

	block, _ := pem.Decode(pbk)
	//enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	/*
		if enc {
			b, err = x509.DecryptPEMBlock(block, nil)
			if err != nil {
				return nil, err
			}
		}
	*/
	key, err := x509.ParsePKCS1PublicKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil

}
