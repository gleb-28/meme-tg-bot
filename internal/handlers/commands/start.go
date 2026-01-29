package commands

import (
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/text"

	"gopkg.in/telebot.v4"
)

func Start(ctx telebot.Context) error {
	fsmManager.FSM.GetFSMForUser(ctx.Chat().ID)
	return ctx.Send(text.Replies.Start)
}
