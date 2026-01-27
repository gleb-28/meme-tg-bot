package config

import (
	"errors"
	"os"
)

type Config struct {
	TgBotToken     string
	LoggerBotToken string
	AdminUsername  string
	ActivationKey  string
	IsDebug        bool
}

func GetConfig() (Config, error) {
	config := Config{}

	config.TgBotToken = os.Getenv("TG_BOT_TOKEN")
	if config.TgBotToken == "" {
		return config, errors.New("environment variable TG_BOT_TOKEN is not set")
	}

	config.LoggerBotToken = os.Getenv("LOGGER_BOT_TOKEN")

	config.AdminUsername = os.Getenv("ADMIN_USERNAME")

	config.ActivationKey = os.Getenv("ACTIVATION_KEY")
	if config.ActivationKey == "" {
		return config, errors.New("environment variable ACTIVATION_KEY is not set")
	}

	config.IsDebug = os.Getenv("IS_DEBUG") == "true"

	return config, nil
}
