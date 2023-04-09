package fpe

import "aiot-service-for-mfp/pkg/kms/algorithms/fpe/ff1"

type Wrapper struct {
	radix   int
	maxTLen int
	key     []byte
	tweak   []byte
}

func (c *Wrapper) SetKey(key []byte) {
	c.key = key
	c.tweak = []byte("longforfpe")
	c.maxTLen = 16
}
func (c *Wrapper) Encrypt(radix int, plaintext string) (string, error) {
	ff, err := ff1.NewCipher(radix, c.maxTLen, c.key, c.tweak)
	if err != nil {
		return "", err
	}

	input, err := ff1.GetTransform().TransformString(plaintext)
	if err != nil {
		return "", err
	}
	ciphertext, err := ff.Encrypt(input)
	if err != nil {
		return "", err
	}

	return ff1.GetTransform().TransformInt(ciphertext)
}
func (c *Wrapper) Decrypt(radix int, ciphertext string) (string, error) {
	ff, err := ff1.NewCipher(radix, c.maxTLen, c.key, c.tweak)
	if err != nil {
		return "", err
	}

	input, err := ff1.GetTransform().TransformString(ciphertext)
	if err != nil {
		return "", err
	}
	plaintext, err := ff.Decrypt(input)
	if err != nil {
		return "", err
	}
	return ff1.GetTransform().TransformInt(plaintext)
}
