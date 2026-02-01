package message

import (
	"context"
	b "memetgbot/internal"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/models"

	"gopkg.in/telebot.v4"
)

func createValidateActivationKey(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID

		if ctx.Message().Text == bot.Config.ActivationKey {
			err := bot.ChatRepo.Upsert(&models.Chat{TelegramID: chatId, Name: ctx.Sender().FirstName})
			if err != nil {
				bot.MustSend(chatId, bot.Replies.Error)
				bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
				return nil
			}
			bot.MustSend(chatId, bot.Replies.EnterKeySuccess)
			bot.MustSend(chatId, bot.Replies.Start)
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			return nil
		}

		bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
		return ctx.Send(bot.Replies.EnterKeyWrong)
	}
}
