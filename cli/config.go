package main

var cfg *Config

type Config struct {
	Verbose      bool
	Location     string
	EnableSingle bool `mapstructure:"enable-single"`
	Mode         ModeType
	Decode       bool
	SinglePath   string `mapstructure:"single-path"`
	KeyPath      string `mapstructure:"key-path"`
	Server       ServerConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type ModeType string

const (
	ModeGenLocal  ModeType = "genlocal"
	ModeRemote    ModeType = "remote"
	ModeLocal     ModeType = "local"
	ModeGenRemote ModeType = "genremote"
)
