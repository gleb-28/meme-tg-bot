package message

import (
	"context"
	b "memetgbot/internal"
	"memetgbot/internal/core/config"
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
			b.SendWithHandlingErr(chatId, text.Replies.Error)
			fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			return nil
		}
		b.SendWithHandlingErr(chatId, text.Replies.EnterKeySuccess)
		b.SendWithHandlingErr(chatId, text.Replies.Start)
		fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
		return nil
	}

	fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
	return ctx.Send(text.Replies.EnterKeyWrong)
}
