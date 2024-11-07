package config

import (
	"fmt"

	"github.com/goloop/env"
)

type Config struct {
	DSN  string
	Host string
	Port string
}

func LoadConfig() (Config, error) {
	config := Config{}
	err := env.Load("config.env")
	if err != nil {
		return config, fmt.Errorf(`{"error":"not possible to load environment variables from a file"}`)
	}
	config.DSN = env.Get("DSN")
	if config.DSN == "" {
		config.DSN = "localhost user=postgres password=postgres dbname=wallet sslmode=disable"
	}
	config.Host = env.Get("HOST")
	if config.Host == "" {
		config.Host = "localhost"
	}
	config.Port = env.Get("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}
	return config, nil
}
