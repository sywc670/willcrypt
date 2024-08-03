//go:build !windows

package utils

import "os"

func GetHomeDir() string {
	return os.Getenv("HOME")
}
