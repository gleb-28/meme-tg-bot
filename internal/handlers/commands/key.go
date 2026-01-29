package commands

import (
	"context"
	"memetgbot/internal/core/logger"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/repo"
	"memetgbot/internal/text"

	"gopkg.in/telebot.v4"
)

func Key(ctx telebot.Context) error {
	chatId := ctx.Chat().ID
	chat, err := repo.Chat.Get(chatId)
	if err != nil {
		logger.Logger.Error(err.Error())
		return ctx.Send(text.Replies.Error)
	}
	if chat.TelegramID == 0 {
		fsm := fsmManager.FSM.GetFSMForUser(chatId)
		err = fsm.Event(context.Background(), fsmManager.AwaitingKeyEvent)
		if err != nil {
			logger.Logger.Error(err.Error())
			return ctx.Send(text.Replies.Error)
		}

		return ctx.Send(text.Replies.EnterKey)
	}

	return ctx.Send(text.Replies.KeyAlreadyEntered)
}
