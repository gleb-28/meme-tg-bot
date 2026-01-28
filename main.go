package main

import (
	c "memetgbot/src/core/config"
	l "memetgbot/src/core/logger"
	"memetgbot/src/db"
	"os"

	g "github.com/joho/godotenv"
	//t "gopkg.in/telebot.v4"
)

func main() {
	if os.Getenv("IS_DOCKERIZED") != "true" {
		err := g.Load()
		if err != nil {
			panic("Error loading .env file (might be missing in non-dockerized env):" + err.Error())
		}
	}

	config, err := c.GetConfig()
	if err != nil {
		panic("Error getting config: " + err.Error())
	}

	logger := l.GetLoggerService(config.IsDebug, l.NewLoggerBot(config.LoggerBotToken, config.AdminID), config.AdminID)

	logger.Error("Hi!")

	d := db.InitDB(&config.Database, logger)

	s := d.Table("chat")
	logger.Info(s.Name())
}
