package config

import (
	"errors"
	"os"
)

type Config struct {
	AppConfig
}

type AppConfig struct {
	Port string
}

func Load() (*Config, error) {
	var appCfg AppConfig
	var err error

	if appCfg.Port, err = getStrValue("APP_PORT", true); err != nil {
		return nil, err
	}

	return &Config{appCfg}, nil
}

func getStrValue(key string, required bool) (string, error) {
	value := os.Getenv(key)
	if !required && value == "" {
		return "", errors.New("varriable is empty")
	}
	return value, nil
}
