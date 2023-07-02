package main

import (
	"os"
	"github.com/BurntSushi/toml"
)

type Config struct {
	URLS              []string
	LinksProDurchlauf int
	AlleExportieren   bool
	TeilExporte       bool
	InfluexUrl        string
	InfluexOrg        string
	InfluexBucket     string
}

type Secrets struct {
	Token string
}

func readConfigFile(filename string) []byte {
	// read file
	b, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return b
}

func getConfig() Config {
	var conf Config
	_, err := toml.Decode(string(readConfigFile("config.toml")), &conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func getSecrets() Secrets {
	var secrets Secrets
	_, err := toml.Decode(string(readConfigFile("secrets.toml")), &secrets)
	if err != nil {
		panic(err)
	}
	return secrets
}
