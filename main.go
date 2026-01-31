package main

import (
	bot "memetgbot/internal"
	"memetgbot/internal/core/logger"
	"memetgbot/internal/db"
	"memetgbot/internal/handlers/commands"
	"memetgbot/internal/handlers/message"
)

func main() {
	var b = bot.Bot

	db.InitDB()

	commands.InitCommands(b)
	message.InitMessagesHandler(b)
	logger.Logger.Info("Bot successfully started!")
	b.Start()
}
