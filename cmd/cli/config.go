package main

var c *Config = NewConfig()

type Config struct {
	Debug          bool
	TargetDir      string
	IsSingle       bool
	UseMode        Mode
	IsDecode       bool
	SingleFilepath string
	StoreKey       string
}

func NewConfig() *Config {
	return &Config{}
}

var mode string

type Mode string

const (
	ModeGenLocal  Mode = "genlocal"
	ModeRemote    Mode = "remote"
	ModeLocal     Mode = "local"
	ModeGenRemote Mode = "genremote"
)
