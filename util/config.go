package util

import (
	"github.com/tkanos/gonfig"
)

// Config stores all configuration of the application.
type Config struct {
	DB_DRIVER string
	DB_SOURCE string
}

// LoadConfig reads config from json file.
func LoadConfig(path string) Config {

	config := Config{}

	if len(path) == 0 {
		path = "./dev_config.json"
	}

	err := gonfig.GetConf(path, &config)
	if err != nil {
		panic(err)
	}

	return config
}
