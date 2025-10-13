package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	Port string `env:"PORT"`
}

func NewConfig() Config {
	return Config{}
}

func (cfg *Config) InitConfig() {
	serverAddress := viper.GetString("SERVER_ADDRESS")
	port := viper.GetString("PORT")
	cfg.ServerAddress = serverAddress
	cfg.Port = port
}