package message

import (
	b "memetgbot/internal"
	"memetgbot/pkg/utils"

	"gopkg.in/telebot.v4"
)

func createHandleMessage(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		textMsg := ctx.Message().Text

		if textMsg != "" && utils.IsURL(textMsg) {
			return createHandleLink(bot)(ctx)
		}

		if forwardChatId, enabled := bot.ForwardModeService.GetForwardChat(chatId); enabled {
			bot.ForwardAnyMessage(ctx, forwardChatId)
			return nil
		}

		return ctx.Send(bot.Replies.CannotChatYet)
	}
}
