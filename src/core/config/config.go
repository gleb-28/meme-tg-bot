package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	TgBotToken     string
	LoggerBotToken string
	AdminID        uint32
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

	adminId, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil {
		return config, errors.New("environment variable ADMIN_ID is not set")
	}
	config.AdminID = uint32(adminId)

	config.ActivationKey = os.Getenv("ACTIVATION_KEY")
	if config.ActivationKey == "" {
		return config, errors.New("environment variable ACTIVATION_KEY is not set")
	}

	config.IsDebug = os.Getenv("IS_DEBUG") == "true"

	return config, nil
}
