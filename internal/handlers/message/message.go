package message

import (
	fsmManager "memetgbot/internal/fsm"

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
		return nil
	case fsmManager.StateAwaitingKey:
		return validateActivationKey(ctx)
	default:
		return nil
	}
}
