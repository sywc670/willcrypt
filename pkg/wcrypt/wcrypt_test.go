package wcrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"

	"github.com/sywc670/willcrypt/internal/config"
)

func TestEncryptAndDecrypt(t *testing.T) {
	priv, _ := rsa.GenerateKey(rand.Reader, config.Bits)

	filename := "testEncryptAndDecrypt.txt"
	testContent := "hello world"

	err := os.WriteFile(filename, []byte(testContent), 0777)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Start encrypt")
	Encrypt(filename, priv)

	if _, err := os.Stat(filename + config.LockedExtension); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filename); err == nil {
		t.Fatalf("%s has not been removed.", filename)
	}

	t.Log("Start decrypt")
	Decrypt(filename+config.LockedExtension, priv)

	if _, err := os.Stat(filename + config.LockedExtension); err == nil {
		t.Fatalf("%s has not been removed.", filename+config.LockedExtension)
	}

	if data, _ := os.ReadFile(filename); string(data) != testContent {
		t.Fatalf("%s data do not match original, got: %s", filename, string(data))
	}
}
