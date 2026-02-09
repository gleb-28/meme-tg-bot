package commands

import (
	bot "memetgbot/internal"
	"memetgbot/internal/middleware/auth"

	"gopkg.in/telebot.v4"
)

var commands = []telebot.Command{
	{Text: "start", Description: "Старт"},
	{Text: "key", Description: "Ввести секретный ключ"},
	{Text: "change_mode", Description: "Поменять режим"},
}

func MustInitCommandsHandler(bot *bot.Bot) {
	err := bot.SetCommands(commands)
	if err != nil {
		panic("Failed to set commands: " + err.Error())
	}

	bot.Handle("/start", createStartHandler(bot), auth.CreateAuthMiddleware(bot))
	bot.Handle("/key", createKeyHandler(bot))
	bot.Handle("/change_mode", createChangeModeHandler(bot), auth.CreateAuthMiddleware(bot))
}
