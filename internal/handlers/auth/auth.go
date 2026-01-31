package auth

import (
	"memetgbot/internal/core/logger"
	"memetgbot/internal/repo"
	"memetgbot/internal/text"

	"gopkg.in/telebot.v4"
)

func WithAuth(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		chat, err := repo.Chat.Get(chatId)
		if err != nil {
			logger.Logger.Error(err.Error())
			return ctx.Send(text.Replies.Error)
		}
		if chat.TelegramID == 0 {
			return nil
		}

		return next(ctx)
	}
}
