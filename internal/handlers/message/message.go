package message

import (
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/handlers/auth"

	"gopkg.in/telebot.v4"
)

func InitMessagesHandler(bot *telebot.Bot) {
	bot.Handle(telebot.OnText, messageHandler)
}

func messageHandler(ctx telebot.Context) error {
	chatId := ctx.Chat().ID
	fsm := fsmManager.FSM.GetFSMForUser(chatId)

	switch fsm.Current() {
	case fsmManager.StateInitial:
		return auth.WithAuth(handleMessage)(ctx)
	case fsmManager.StateAwaitingKey:
		return validateActivationKey(ctx)
	case fsmManager.StateProcessingLink:
		return auth.WithAuth(handleProcessingLink)(ctx)
	default:
		return nil
	}
}
