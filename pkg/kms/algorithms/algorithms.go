package algorithms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"math/rand"
	"time"
)

func ShaHmac1(source, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(source))
	signedBytes := h.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Seed(time.Now().Unix())
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(s)
}
