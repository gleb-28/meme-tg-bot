package message

import (
	"context"
	"memetgbot/internal/core/config"
	"memetgbot/internal/core/logger"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/repo"
	"memetgbot/internal/text"

	"gopkg.in/telebot.v4"
)

func validateActivationKey(ctx telebot.Context) error {
	chatId := ctx.Chat().ID

	if ctx.Message().Text == config.Config.ActivationKey {
		err := repo.Chat.Add(chatId)
		if err != nil {
			ctx.Send(text.Replies.Error)
			return setInitialState(ctx, chatId)
		}
		ctx.Send(text.Replies.EnterKeySuccess)
		return setInitialState(ctx, chatId)
	} else {
		ctx.Send(text.Replies.EnterKeyWrong)
		return setInitialState(ctx, chatId)
	}
}

func setInitialState(ctx telebot.Context, chatId int64) error {
	fsm := fsmManager.FSM.GetFSMForUser(chatId)
	err := fsm.Event(context.Background(), fsmManager.InitialEvent)
	if err != nil {
		logger.Logger.Error(err.Error())
		return ctx.Send(text.Replies.Error)
	}
	return nil
}
