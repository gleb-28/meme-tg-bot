package message

import (
	b "memetgbot/internal"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/middleware/auth"

	"gopkg.in/telebot.v4"
)

func MustInitMessagesHandler(bot *b.Bot) {
	bot.Handle(telebot.OnText, createMessageHandler(bot))
	bot.Handle(telebot.OnAudio, createMessageHandler(bot))
	bot.Handle(telebot.OnDocument, createMessageHandler(bot))
	bot.Handle(telebot.OnPhoto, createMessageHandler(bot))
	bot.Handle(telebot.OnVideo, createMessageHandler(bot))
	bot.Handle(telebot.OnVideoNote, createMessageHandler(bot))
	bot.Handle(telebot.OnAnimation, createMessageHandler(bot))
	bot.Handle(telebot.OnSticker, createMessageHandler(bot))
	bot.Handle(telebot.OnVoice, createMessageHandler(bot))
}

func createMessageHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		userFsm := bot.Fsm.GetFSMForUser(chatId)

		switch userFsm.Current() {
		case fsmManager.StateInitial:
			go func(bot *b.Bot, ctx telebot.Context) {
				_ = auth.CreateAuthMiddleware(bot)(createHandleMessage(bot))(bot.NewContext(ctx.Update()))
			}(bot, ctx)
			return nil
		case fsmManager.StateAwaitingKey:
			return createValidateActivationKey(bot)(ctx)
		case fsmManager.StateAwaitingForwardChat:
			return auth.CreateAuthMiddleware(bot)(createValidateForwardChatHandler(bot))(ctx)
		default:
			return nil
		}
	}
}
