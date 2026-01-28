package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TgBotToken     string
	LoggerBotToken string
	AdminID        uint32
	ActivationKey  string
	IsDebug        bool
	Database       DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
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

	config.Database = DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
	if config.Database.Host == "" {
		return config, errors.New("environment variable DB_HOST is not set")
	}
	if config.Database.Port == "" {
		return config, errors.New("environment variable DB_PORT is not set")
	}
	if config.Database.User == "" {
		return config, errors.New("environment variable DB_USER is not set")
	}
	if config.Database.Password == "" {
		return config, errors.New("environment variable DB_PASSWORD is not set")
	}
	if config.Database.Name == "" {
		return config, errors.New("environment variable DB_NAME is not set")
	}

	return config, nil
}

func (dc *DatabaseConfig) ToDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dc.Host, dc.User, dc.Password, dc.Name, dc.Port)
}
