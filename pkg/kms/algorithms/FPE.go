package algorithms

import "aiot-service-for-mfp/pkg/kms/algorithms/fpe"

type FPE struct {
}

func (c *FPE) Encrypt(key []byte, radix int, str string) (string, error) {
	if len(str) < 2 {
		return str, nil
	}
	fw := fpe.Wrapper{}
	fw.SetKey(key)
	return fw.Encrypt(radix, str)
}
func (c *FPE) Decrypt(key []byte, radix int, str string) (string, error) {
	if len(str) < 2 {
		return str, nil
	}
	fw := fpe.Wrapper{}
	fw.SetKey(key)
	return fw.Decrypt(radix, str)
}
