package main

import (
	"github.com/docopt/docopt.go"
	"log"
)

var (
	githash string = ""
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	usage := `Preview Deploy

Usage: preview_deploy [--help --version --config=<file>]
       
Options:
  --help           Show this screen.
  --version        Show version.
  --config=<file>  The configuration file to use.`

	arguments, _ := docopt.Parse(usage, nil, true, version(), false)
	configPath := GetConfigString(arguments, "--config")
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading config", err)
	}
	app, err := NewApp(config)
	if err != nil {
		log.Fatal("Error creating app", err)
	}
	app.Start()
}

func version() string {
	previewVersion := "1.2.0"
	if len(githash) > 0 {
		return previewVersion + "+" + githash
	}
	return previewVersion
}

func GetConfigString(arguments map[string]interface{}, key string) string {
	configPath, hasConfigPath := arguments[key]
	if hasConfigPath {
		value, ok := configPath.(string)
		if ok {
			return value
		}
	}
	return ""
}
