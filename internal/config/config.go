package config

import (
	"os"
)

type Config struct {
	SrvHost    string
	SrvPort    string
	DbDatabase string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
}

func env(key string, or string) string {
	val := os.Getenv(key)

	if val == "" {
		return or
	}

	return val
}

func New() *Config {
	return &Config{
		SrvHost:    env("SRV_HOST", "localhost"),
		SrvPort:    env("SRV_PORT", "8000"),
		DbDatabase: env("DB_DATABASE", "postgres"),
		DbHost:     env("DB_HOST", "localhost"),
		DbPort:     env("DB_PORT", "5432"),
		DbUser:     env("DB_USER", "postgres"),
		DbPassword: env("DB_PASSWORD", "postgres"),
	}
}
