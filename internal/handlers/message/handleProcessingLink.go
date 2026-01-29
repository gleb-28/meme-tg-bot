package message

import "gopkg.in/telebot.v4"

func handleProcessingLink(ctx telebot.Context) error {
	return ctx.Send("Погоди, не все сразу")
}
