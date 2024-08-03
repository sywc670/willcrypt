package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
)

func GenerateID() string {
	b := make([]byte, 32)
	rand.Read(b)

	hash := sha256.New()

	return hex.EncodeToString(hash.Sum(b))
}

func Stringify(priv *rsa.PrivateKey) string {
	b := x509.MarshalPKCS1PrivateKey(priv)
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: b,
	}

	return string(pem.EncodeToMemory(&block))
}
