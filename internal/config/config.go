package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	Port string `mapstructure:"PORT"`
}

func NewConfig() Config {
	return Config{}
}

func (cfg *Config) InitConfig(configPath string) error {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	return nil
}