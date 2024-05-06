package config

import (
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DatabaseConfig DatabaseConfig
	ServerPort     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.DatabaseConfig.Host = os.Getenv("DB_HOST")
	cfg.DatabaseConfig.Port = os.Getenv("DB_PORT")
	cfg.DatabaseConfig.User = os.Getenv("DB_USER")
	cfg.DatabaseConfig.Password = os.Getenv("DB_PASSWORD")
	cfg.DatabaseConfig.Name = os.Getenv("DB_NAME")
	cfg.ServerPort = os.Getenv("PORT")
	return cfg, nil
}
