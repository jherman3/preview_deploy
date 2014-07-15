package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type appConfig struct {
	Source       string      `json:"-"`
	Environments []string    `json:"environments"`
	Nodes        []*chefNode `json:"nodes"`
}

type chefNode struct {
	Hostname string `json:"hostname"`
	Url      string `json:"url"`
}

func loadConfig(path string) (*appConfig, error) {
	configPath := determineConfigPath(path)
	if configPath == "" {
		log.Fatal("No configuration file could be found")
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	return newConfig(data)
}

func newConfig(data []byte) (*appConfig, error) {
	var config appConfig
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	config.Source = string(data)
	return &config, nil
}

func determineConfigPath(givenPath string) string {
	wd, _ := os.Getwd()
	paths := []string{
		givenPath,
		filepath.Join(wd, "preview_deploy.config"),
		"/etc/preview_deploy.config",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
