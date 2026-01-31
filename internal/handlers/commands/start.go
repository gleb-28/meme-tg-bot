package commands

import (
	"memetgbot/internal/text"

	"gopkg.in/telebot.v4"
)

func Start(ctx telebot.Context) error {
	return ctx.Send(text.Replies.Start)
}
