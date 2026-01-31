package message

import (
	b "memetgbot/internal"
	"memetgbot/pkg/utils"

	"gopkg.in/telebot.v4"
)

func createHandleMessage(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		if utils.IsURL(ctx.Message().Text) {
			return createHandleLink(bot)(ctx)
		}
		return ctx.Send("Я пока не умею общаться(") // TODO
	}
}
