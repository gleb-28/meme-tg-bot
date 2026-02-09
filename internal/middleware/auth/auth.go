package auth

import (
	b "memetgbot/internal"

	"gopkg.in/telebot.v4"
)

func CreateAuthMiddleware(bot *b.Bot) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(ctx telebot.Context) error {
			chatId := ctx.Chat().ID
			chat, err := bot.GetChatCached(chatId)
			if err != nil {
				bot.Logger.Error(err.Error())
				return ctx.Send(bot.Replies.Error)
			}
			if chat.TelegramID == 0 {
				return nil
			}

			return next(ctx)
		}
	}
}
