package commands

import (
	b "memetgbot/internal"
	"memetgbot/internal/handlers/keyboard"

	"gopkg.in/telebot.v4"
)

func createChangeModeHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		isEnabled := bot.ForwardModeService.IsForwardEnabled(chatId)
		bot.MustSend(chatId, bot.Replies.SelectTheMode, keyboard.ForwardModeKeyboard(isEnabled, bot.Replies))
		return nil
	}
}
