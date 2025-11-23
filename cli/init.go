package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	initConf()
	check()
}

func initConf() {
	configFile := pflag.StringP("config", "c", filepath.Join(getExeLoc(), "config.yml"), "default relative to exe")
	pflag.String("location", "testground", "relative to workdir")
	verbose := pflag.BoolP("verbose", "v", false, "")
	pflag.StringP("mode", "m", "genlocal", "")
	pflag.StringP("single-path", "f", "", "must enable single")
	pflag.Bool("enable-single", false, "")
	pflag.BoolP("decode", "d", false, "")
	pflag.StringP("key-path", "k", filepath.Join(getExeLoc(), "priv.key"), "default relative to exe")
	pflag.Int("server.port", 8080, "server port")
	pflag.String("server.host", "localhost", "server host")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetConfigFile(*configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if *verbose {
			fmt.Fprintf(os.Stderr, "Read config fail: %v\n", err)
		}
	}

	viper.SetEnvPrefix("wcrypt")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Load config fail: %v\n", err)
	}

}

func check() {
	if cfg.EnableSingle && cfg.SinglePath == "" {
		debugf("Args conflict.\n")
		os.Exit(1)
	}

	if cfg.Mode == ModeGenLocal && cfg.Decode {
		debugf("Args conflict.\n")
		debugf("Change Mode to %s.\n", ModeLocal)
		cfg.Mode = ModeLocal
	}

	if cfg.Mode == ModeGenRemote && cfg.Decode {
		debugf("Args conflict.\n")
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

	// Overwrite loction with args if specified.
	if len(pflag.Args()) >= 1 {
		cfg.Location = pflag.Arg(0)
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

	// if key is found on the keypath, can't use gen mode to avoid overwrite.
	if cfg.Mode == ModeGenLocal || cfg.Mode == ModeGenRemote {
		_, err = os.Stat(cfg.KeyPath)
		if err == nil {
			fmt.Fprintf(os.Stderr, "%s exists with gen mode set is not ok.\n", cfg.KeyPath)
			os.Exit(1)
		}
	}
}

func getExeLoc() string {
	s, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(s)
}
