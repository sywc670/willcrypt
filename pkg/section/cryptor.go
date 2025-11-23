package section

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"io"
	"os"
	"strings"

	"github.com/sywc670/willcrypt/internal/config"
)

type Cryptor struct {
	priv    *rsa.PrivateKey
	bufSize int
}

func NewCryptor(priv *rsa.PrivateKey) *Cryptor {
	return &Cryptor{
		priv:    priv,
		bufSize: 16 * 1024 * 1024, // 默认 16MB
	}
}

// -----------------------------
// 加密
// -----------------------------
func (c *Cryptor) Encrypt(path string) (err error) {
	{
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		tmpPath := path + config.LockedExtension + ".tmp"
		out, err := os.Create(tmpPath)
		if err != nil {
			return err
		}
		defer func() {
			out.Close()
			if err != nil {
				os.Remove(tmpPath)
			}
		}()

		// 1. 生成 AES Key/IV
		key, iv, err := c.genKeyIV()
		if err != nil {
			return err
		}

		// 2. RSA 加密 header
		encHeader, err := c.encryptHeader(key, iv)
		if err != nil {
			return err
		}

		if _, err = out.Write(encHeader); err != nil {
			return err
		}

		// 3. AES 流加密正文
		if err = c.streamEncrypt(in, out, key, iv); err != nil {
			return err
		}

		if err = out.Sync(); err != nil {
			return err
		}
		out.Close()

		finalPath := path + config.LockedExtension
		if err = os.Rename(tmpPath, finalPath); err != nil {
			return err
		}
	}

	_ = os.Remove(path)
	return nil
}

// -----------------------------
// 解密
// -----------------------------
func (c *Cryptor) Decrypt(path string) (err error) {
	{
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		tmpPath := strings.TrimSuffix(path, config.LockedExtension) + ".tmp"
		out, err := os.Create(tmpPath)
		if err != nil {
			return err
		}
		defer func() {
			out.Close()
			if err != nil {
				os.Remove(tmpPath)
			}
		}()

		// 1. 读取并 RSA 解密 header
		key, iv, err := c.decryptHeader(in)
		if err != nil {
			return err
		}

		// 2. AES 解密正文
		if err = c.streamDecrypt(in, out, key, iv); err != nil {
			return err
		}

		if err = out.Sync(); err != nil {
			return err
		}
		out.Close()

		finalPath := strings.TrimSuffix(path, config.LockedExtension)
		if err = os.Rename(tmpPath, finalPath); err != nil {
			return err
		}

	}
	_ = os.Remove(path)
	return nil
}

// -----------------------------
// 内部方法：生成 key / iv
// -----------------------------
func (c *Cryptor) genKeyIV() ([]byte, []byte, error) {
	key := make([]byte, config.KeySize)
	iv := make([]byte, aes.BlockSize)

	if _, err := rand.Read(key); err != nil {
		return nil, nil, err
	}
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}
	return key, iv, nil
}

// -----------------------------
// 内部方法：加密 header
// -----------------------------
func (c *Cryptor) encryptHeader(key, iv []byte) ([]byte, error) {
	header := append(append([]byte{}, key...), iv...)
	pub := c.priv.PublicKey
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &pub, header, []byte(""))
}

// -----------------------------
// 内部方法：解密 header
// -----------------------------
func (c *Cryptor) decryptHeader(r io.Reader) ([]byte, []byte, error) {
	rsaSize := c.priv.PublicKey.Size()
	encHeader := make([]byte, rsaSize)

	if _, err := io.ReadFull(r, encHeader); err != nil {
		return nil, nil, err
	}

	header, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.priv, encHeader, []byte(""))
	if err != nil {
		return nil, nil, err
	}

	key := header[:config.KeySize]
	iv := header[config.KeySize : config.KeySize+aes.BlockSize]
	return key, iv, nil
}

// -----------------------------
// 内部方法：流式 AES 加密
// -----------------------------
func (c *Cryptor) streamEncrypt(in io.Reader, out io.Writer, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	buf := make([]byte, c.bufSize)

	for {
		n, err := in.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				return wErr
			}
		}

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// -----------------------------
// 内部方法：流式 AES 解密
// -----------------------------
func (c *Cryptor) streamDecrypt(in io.Reader, out io.Writer, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	buf := make([]byte, c.bufSize)

	for {
		n, err := in.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, wErr := out.Write(buf[:n]); wErr != nil {
				return wErr
			}
		}

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
