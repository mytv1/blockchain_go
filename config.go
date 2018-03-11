package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	defaultConfigPath = "config.json"
)

// Config contains program's configuration information
type Config struct {
	Nw      Network        `json:"network"`
	SWallet StorableWallet `json:"wallet"`
}

var config *Config

func getConfig() *Config {
	return config
}

func initConfig(configPathCLI string) *Config {
	var configPath string
	if configPathCLI != "" {
		configPath = configPathCLI
	} else {
		configPath = defaultConfigPath
	}

	config = importConfig(configPath)
	wallet = config.SWallet.toWallet()
	setWallet(wallet)
	return config
}

func importConfig(filePath string) *Config {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}

	config := Config{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}
	return &config
}

func (config *Config) exportConfig(filePath string) {
	prettyMarshal, e := json.MarshalIndent(config, "", "  ")
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}

	e = ioutil.WriteFile(filePath, prettyMarshal, 0644)
	if e != nil {
		Error.Println(e.Error())
		os.Exit(1)
	}
}
