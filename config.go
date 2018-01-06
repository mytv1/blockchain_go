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
	Nw Network `json:"network"`
}

var config Config

func getConfig() Config {
	return config
}

func initConfig(configPathCLI string) {
	var configPath string
	if configPathCLI != "" {
		configPath = configPathCLI
	} else {
		configPath = defaultConfigPath
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
