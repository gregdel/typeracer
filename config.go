package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the config
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ConfigFromFile returns the configuration from a file
func ConfigFromFile(filepath string) (*Config, error) {
	file, err := os.Open("config.yml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	return config, yaml.NewDecoder(file).Decode(config)
}
