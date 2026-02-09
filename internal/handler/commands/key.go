package commands

import (
	"context"
	b "memetgbot/internal"
	fsmManager "memetgbot/internal/fsm"

	"gopkg.in/telebot.v4"
)

func createKeyHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		chat, err := bot.GetChatCached(chatId)
		if err != nil {
			bot.Logger.Error(err.Error())
			return ctx.Send(bot.Replies.Error)
		}
		if chat.TelegramID == 0 {
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.AwaitingKeyEvent)
			return ctx.Send(bot.Replies.EnterKey)
		}

		return ctx.Send(bot.Replies.KeyAlreadyEntered)
	}
}
