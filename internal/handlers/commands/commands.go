package commands

import (
	"gopkg.in/telebot.v4"
)

func InitCommands(bot *telebot.Bot) {
	bot.Handle("/start", Start)
}
