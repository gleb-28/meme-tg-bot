package commands

import (
	"memetgbot/internal/handlers/auth"

	"gopkg.in/telebot.v4"
)

var commands = []telebot.Command{
	{Text: "start", Description: "Старт"},
	{Text: "key", Description: "Ввести секретный ключ"},
}

func InitCommands(bot *telebot.Bot) {
	err := bot.SetCommands(commands)
	if err != nil {
		panic("Failed to set commands: " + err.Error())
	}

	bot.Handle("/start", Start, auth.WithAuth)
	bot.Handle("/key", Key)
}
