package config

import (
	"log"
	"os"
	"strconv"

	g "github.com/joho/godotenv"
)

type AppConfig struct {
	TgBotToken     string
	LoggerBotToken string
	AdminID        int64
	ActivationKey  string
	IsDebug        bool
	Database       DatabaseConfig
}

type DatabaseConfig struct {
	FileName string
}

func MustConfig() *AppConfig {
	err := g.Load()
	if err != nil {
		log.Fatal("Error loading .env file (might be missing in non-dockerized env):", err.Error())
	}

	config := &AppConfig{}

	config.TgBotToken = os.Getenv("TG_BOT_TOKEN")
	if config.TgBotToken == "" {
		log.Fatal("environment variable TG_BOT_TOKEN is not set")
	}

	config.LoggerBotToken = os.Getenv("LOGGER_BOT_TOKEN")

	adminId, err := strconv.Atoi(os.Getenv("ADMIN_ID"))
	if err != nil {
		log.Fatal("environment variable ADMIN_ID is not set")
	}
	config.AdminID = int64(adminId)

	config.ActivationKey = os.Getenv("ACTIVATION_KEY")
	if config.ActivationKey == "" {
		log.Fatal("environment variable ACTIVATION_KEY is not set")
	}

	config.IsDebug = os.Getenv("IS_DEBUG") == "true"

	config.Database = DatabaseConfig{
		FileName: os.Getenv("DB_FILE_NAME"),
	}
	if config.Database.FileName == "" {
		log.Fatal("environment variable DB_FILE_NAME is not set")
	}

	return config
}
