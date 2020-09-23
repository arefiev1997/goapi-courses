package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// Config base config
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig server config
type ServerConfig struct {
	Port int
}

// DatabaseConfig db config
type DatabaseConfig struct {
	ConnectionString string `toml:"connection_string"`
}

// ReadConfig read config from path
func ReadConfig(path string) (*Config, error) {
	var conf Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file, %v", err)
	}
	if _, err := toml.Decode(string(data), &conf); err != nil {
		// handle error
		return nil, fmt.Errorf("Failed to decode config %v", err)
	}
	return &conf, nil
}
