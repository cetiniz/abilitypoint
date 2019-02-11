package config

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL  string `env:"BASE_URL"`
	UserName string `env:"USER_NAME"`
	UserPass string `env:"USER_PASS"`
}

func (cfg *Config) LoadEnvironmentVariables() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("File .env not found, reading configuration from ENV")
	}
	if err := env.Parse(cfg); err != nil {
		fmt.Println("Failed to parse ENV")
	}
}

func (cfg *Config) ReturnEnvironmentVariables() (string, string, string) {
	return cfg.BaseURL, cfg.UserName, cfg.UserPass
}
