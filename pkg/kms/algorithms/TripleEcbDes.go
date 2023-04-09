package algorithms

import (
	"aiot-service-for-mfp/pkg/kms/algorithms/tripledes"
	"errors"
)

type TripleEcbDes struct {
}

func NewTripleEcbDes() *TripleEcbDes {
	return &TripleEcbDes{}
}
func (c *TripleEcbDes) Encrypt(privateKey, plaintext []byte) ([]byte, error) {
	if len(privateKey) != 16 {
		return nil, errors.New("invalid key length")
	}

	key1 := privateKey[:8]
	key := append(privateKey, key1...)
	return tripledes.TripleEcbDesEncrypt(plaintext, key)
}
func (c *TripleEcbDes) Decrypt(privateKey, in []byte) ([]byte, error) {
	if len(privateKey) != 16 {
		return nil, errors.New("invalid key length")
	}

	key1 := privateKey[:8]
	key := append(privateKey, key1...)
	return tripledes.TripleEcbDesDecrypt(in, key)
}
