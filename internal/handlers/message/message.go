package message

import (
	b "memetgbot/internal"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/handlers/auth"

	"gopkg.in/telebot.v4"
)

func MustInitMessagesHandler(bot *b.Bot) {
	bot.Handle(telebot.OnText, createMessageHandler(bot))
}

func createMessageHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		userFsm := bot.Fsm.GetFSMForUser(chatId)

		switch userFsm.Current() {
		case fsmManager.StateInitial:
			return auth.CreateAuthMiddleware(bot)(createHandleMessage(bot))(ctx)
		case fsmManager.StateAwaitingKey:
			return createValidateActivationKey(bot)(ctx)
		case fsmManager.StateProcessingLink:
			return auth.CreateAuthMiddleware(bot)(handleProcessingLink)(ctx)
		default:
			return nil
		}
	}
}
