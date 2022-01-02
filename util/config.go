package util

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

// Config stores all configuration of the application.
type Config struct {
	DB_DRIVER   string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SCHEMA   string
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

// GetDBSource creates the db source from a config file.
func GetDBSource(config Config) string {

	if len(config.DB_SCHEMA) == 0 {
		config.DB_SCHEMA = "public"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=%s",
		config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_SCHEMA)
}
