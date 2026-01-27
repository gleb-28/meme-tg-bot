package main

import (
	c "memetgbot/src/core/config"
	l "memetgbot/src/core/logger"

	g "github.com/joho/godotenv"
	//t "gopkg.in/telebot.v4"
)

func main() {
	err := g.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config, err := c.GetConfig()
	if err != nil {
		panic("Error getting config: " + err.Error())
	}

	logger := l.GetLoggerService(config.IsDebug, l.NewLoggerBot(config.LoggerBotToken, config.AdminID), config.AdminID)

	logger.Error("Hi!")

}
