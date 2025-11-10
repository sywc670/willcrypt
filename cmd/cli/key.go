package main

import (
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sywc670/willcrypt/internal/utils"
)

func getPrivKey() (priv *rsa.PrivateKey, err error) {
	switch c.UseMode {
	case ModeGenLocal, ModeGenRemote:
		fmt.Println("Generating PrivateKey")
		priv = utils.GenerateKey()
	case ModeLocal:
		// TODO: only support one key.
		key, err := readKeyFromFile()
		if err != nil {
			panic(err)
		}
		priv, err := utils.DecodeKey([]byte(key))
		if err != nil {
			panic(err)
		}
		return priv, err
	case ModeRemote:
		idFile, err := os.Open("id.txt")

		if err == nil {
			content, err := io.ReadAll(idFile)
			idFile.Close()

			if err != nil {
				return nil, err
			}

			id := strings.Split(string(content), "\r\n")[1]

			debugf("id: %s\n", id)

			priv = utils.GetKey(id)
		} else {
			return nil, err
		}
	}
	return
}
func readKeyFromFile() (key string, err error) {
	_, err = os.Stat(c.StoreKey)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(c.StoreKey, "not exists")
			return
		} else {
			panic(err)
		}
	}

	bs, err := os.ReadFile(c.StoreKey)
	if err != nil {
		panic(err)
	}

	return utils.DecodeBase64(string(bs)), nil

}

func storeOrUpload(priv *rsa.PrivateKey) {
	switch c.UseMode {
	case ModeGenLocal:
		file, err := os.OpenFile(c.StoreKey, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(utils.EncodeBase64(utils.Stringify(priv)))
		if err != nil {
			panic(err)
		}
	case ModeGenRemote:
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
}
