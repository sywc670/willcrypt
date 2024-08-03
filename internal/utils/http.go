package utils

import (
	"crypto/rsa"
	"io"
	"net/http"
	"net/url"

	"github.com/sywc670/willcrypt/internal/config"
)

func PostKey(priv *rsa.PrivateKey, id string) error {
	key := Stringify(priv)

	// MORE: tls
	_, err := http.PostForm(config.UploadEndpoint, url.Values{
		"key": {key},
		"id":  {id},
	})

	return err
}

func GetKey(id string) *rsa.PrivateKey {
	resp, err := http.PostForm(config.RetrieveEndpoint, url.Values{
		"id": []string{id},
	})
	if err != nil {
		panic(err)
	}

	key := make([]byte, resp.ContentLength)
	io.ReadFull(resp.Body, key)

	priv, err := DecodeKey(key)
	if err != nil {
		panic(err)
	}
	return priv
}
