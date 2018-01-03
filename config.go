package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	DEFAULT_CONFIG_PATH = "config.json"
)

type Config struct {
	Nw Network `json:"network"`
}

var config Config

func GetConfig() Config {
	return config
}

func InitConfig(configPathCLI string) {
	var configPath string
	if configPathCLI != "" {
		configPath = configPathCLI
	} else {
		configPath = DEFAULT_CONFIG_PATH
	}

	config = importConfig(configPath)
}

func importConfig(filePath string) Config {
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
	return config
}
