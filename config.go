package main

import (
	"encoding/json"
	"os"
)

// Config represents basic configurations for the server
type Config struct {
	Server   string
	Port     string
	Channels []string
	Password string
}

// NewConfig creates a new config object from a json file
func NewConfig(path string) *Config {
	file, err := os.Open(path)
	defer file.Close()
	check(err)

	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	check(err)

	return &conf
}
