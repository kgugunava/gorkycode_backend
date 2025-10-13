package config

import (
	// "github.com/caarlos0/env/v11"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	Port string `env:"PORT"`
	DbUser string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbAddress string `env:"DB_ADDRESS"`
	DbPort string `env:"DB_PORT"`
	SslMode string `env:"SSL_MODE"`
	DbName string `env:"DB_NAME"`
}

func NewConfig() Config {
	return Config{}
}

func (cfg *Config) InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	cfg.ServerAddress = os.Getenv("SERVER_ADDRESS")
	cfg.Port = os.Getenv("SERVER_PORT")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.DbPassword = os.Getenv("DB_PASSWORD")
	cfg.DbAddress = os.Getenv("DB_ADDRESS")
	cfg.DbPort = os.Getenv("DB_PORT")
	cfg.SslMode = os.Getenv("SSL_MODE")
	cfg.DbName = os.Getenv("DB_NAME")
	return nil
}