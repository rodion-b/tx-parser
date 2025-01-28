package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTP_ADDR string
}

// Reads config from environment.
func Read() (*Config, error) {
	config := Config{}
	HTTP_ADDR, exists := os.LookupEnv("HTTP_ADDR")
	if exists {
		config.HTTP_ADDR = HTTP_ADDR
	} else {
		return nil, errors.New("HTTP_ADDR is not set")
	}

	return &config, nil
}
