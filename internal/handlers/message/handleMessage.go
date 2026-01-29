package message

import (
	"memetgbot/pkg/utils"

	"gopkg.in/telebot.v4"
)

func handleMessage(ctx telebot.Context) error {
	if utils.IsURL(ctx.Message().Text) {
		return handleLink(ctx)
	}

	return ctx.Send("Я пока не умею общаться(")
}
