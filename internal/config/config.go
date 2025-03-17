package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig   `mapstructure:"server"`
	DatabaseConfig `mapstructure:"database"`
	APIConfig      `mapstructure:"financial_api"`
}

type ServerConfig struct {
	ServerHost      string        `mapstructure:"host"`
	ServerPort      int           `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	DBHost     string `mapstructure:"host"`
	DBPort     int    `mapstructure:"port"`
	DBUser     string `mapstructure:"user"`
	DBPassword string `mapstructure:"password"`
	DBName     string `mapstructure:"name"`
}

type APIConfig struct {
	AlphaVantage `mapstructure:"alphavantage"`
}

type AlphaVantage struct {
	Url    string `mapstructure:"base_url"`
	ApiKey string `mapstructure:"api_key"`
}

func LoadConfig() (*Config, error) {
	configPath, err := env("CONFIG_PATH")

	if err != nil {
		log.Error().Caller().Err(err).Msg("Not able to find a path to the configuration file")
		return nil, err
	}

	if len(configPath) == 0 {
		log.Error().Caller().Msg("The path to the configuration file is empty")
		return nil, err
	}

	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		err := fmt.Errorf("[LoadConfig] error reading config: %w", err)
		log.Error().Caller().Err(err).Msg("viper failed to read config data")
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error().Caller().Err(err).Msg("Failed to unmarshaling config data")
		return nil, err
	}

	return &cfg, nil
}

func env(key string) (string, error) {
	if value, exist := os.LookupEnv(key); exist {
		return value, nil
	}
	return "", fmt.Errorf("ENV %s does not exist", key)
}
