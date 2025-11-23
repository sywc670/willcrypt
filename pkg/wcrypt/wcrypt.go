package wcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sywc670/willcrypt/internal/config"
)

func EncryptBySection(path string, priv *rsa.PrivateKey) (err error) {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if in != nil {
			in.Close()
			in = nil
		}
	}()

	// 临时文件（写完后再原子替换）
	tmpPath := path + config.LockedExtension + ".tmp"

	out, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	// 若失败，确保临时文件删除
	defer func() {
		out.Close()
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	// --- 生成 AES key + IV ---
	key := make([]byte, config.KeySize)
	iv := make([]byte, aes.BlockSize)
	rand.Read(key)
	rand.Read(iv)

	// --- RSA 加密 header (key + iv) ---
	header := append([]byte{}, key...)
	header = append(header, iv...)

	pub := priv.PublicKey
	encHeader, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &pub, header, []byte(""))
	if err != nil {
		return err
	}

	// 写入 header
	if _, err = out.Write(encHeader); err != nil {
		return err
	}

	// --- AES 流加密 ---
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)

	buf := make([]byte, 100*1024*1024) // 100MB buffer

	for {
		n, readErr := in.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, err = out.Write(buf[:n]); err != nil {
				return err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	// 确保写盘（断电时也尽可能保证完整）
	if err = out.Sync(); err != nil {
		return err
	}

	out.Close()

	// --- 原子替换 ---
	finalPath := path + config.LockedExtension
	if err = os.Rename(tmpPath, finalPath); err != nil {
		return err
	}

	if in != nil {
		in.Close()
		in = nil
	}

	// 删除原文件（失败也不会影响最终加密结果）
	if rmErr := os.Remove(path); rmErr != nil {
		log.Printf("Remove original file failed: %v\n", rmErr)
	}

	return nil
}

func DecryptBySection(path string, priv *rsa.PrivateKey) (err error) {
	// 打开加密文件
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if in != nil {
			in.Close()
			in = nil
		}
	}()

	// 临时文件（写完后原子替换）
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

	// --- 先读取被 RSA 加密的 header ---
	rsaHeaderSize := priv.PublicKey.Size() // RSA ciphertext size = key length
	encHeader := make([]byte, rsaHeaderSize)

	n, err := io.ReadFull(in, encHeader)
	if err != nil || n != rsaHeaderSize {
		return fmt.Errorf("读取 header 失败: %v", err)
	}

	// --- RSA 解密得到 key + iv ---
	header, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, encHeader, []byte(""))
	if err != nil {
		return fmt.Errorf("RSA 解密 header 失败: %v", err)
	}

	key := header[:config.KeySize]
	iv := header[config.KeySize : config.KeySize+aes.BlockSize]

	// --- AES-CFB 解密器 ---
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	// --- 流式解密正文 ---
	buf := make([]byte, 100*1024*1024) // 100MB buffer
	for {
		n, readErr := in.Read(buf)
		if n > 0 {
			// 就地解密
			stream.XORKeyStream(buf[:n], buf[:n])

			// 写入解密结果
			if _, err := out.Write(buf[:n]); err != nil {
				return err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	// 保证写盘
	if err = out.Sync(); err != nil {
		return err
	}

	out.Close()

	// --- 原子替换 ---
	finalPath := strings.TrimSuffix(path, config.LockedExtension)
	if err = os.Rename(tmpPath, finalPath); err != nil {
		return err
	}

	if in != nil {
		in.Close()
		in = nil
	}

	// 删除加密文件
	if rmErr := os.Remove(path); rmErr != nil {
		log.Printf("删除加密文件失败: %v\n", rmErr)
	}

	return nil
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
