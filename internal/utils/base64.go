package utils

import "encoding/base64"

func DecodeBase64(s string) string {
	b, _ := base64.StdEncoding.DecodeString(s)
	return string(b)
}

func EncodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
