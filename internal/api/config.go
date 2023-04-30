package api

import (
	"iContext/repository/postgres"
)

// General instance for API server of REST application
type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LoggerLevel string `toml:"logger_level"`
	Storage     *postgres.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		LoggerLevel: "debug",
		Storage:     postgres.NewConfig(),
	}
}
