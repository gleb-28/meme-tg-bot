package keyboard

import (
	b "memetgbot/internal"
	"memetgbot/internal/core/constants"
	"memetgbot/internal/handlers/auth"
	"strings"

	"gopkg.in/telebot.v4"
)

func MustInitKeyboardHandler(bot *b.Bot) {
	bot.Handle(telebot.OnCallback, createKeyboardHandlers(bot), auth.CreateAuthMiddleware(bot))
}

func createKeyboardHandlers(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		cb := ctx.Callback()
		switch strings.TrimSpace(cb.Data) {
		case constants.EnableForwardModeButtonCallbackQuery:
			return createForwardModeEnabledHandler(bot)(ctx)
		case constants.DisableForwardModeButtonCallbackQuery:
			return createForwardModeDisabledHandler(bot)(ctx)
		case constants.UsePrevForwardChatButtonCallbackQuery:
			return createUsePrevForwardChatHandler(bot)(ctx)
		case constants.ChangeForwardChatButtonCallbackQuery:
			return createChangeForwardChatHandler(bot)(ctx)
		default:
			return nil
		}
	}
}
