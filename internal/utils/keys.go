package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/sywc670/willcrypt/internal/config"
)

func GenerateKey() *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, config.Bits)
	if err != nil {
		panic(err)
	}
	return priv
}

func DecodeKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
