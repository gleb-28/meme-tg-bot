package message

import (
	"fmt"
	b "memetgbot/internal"
	"memetgbot/pkg/utils"

	"gopkg.in/telebot.v4"
)

func createHandleMessage(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID

		if utils.IsURL(ctx.Message().Text) {
			return createHandleLink(bot)(ctx)
		}

		if forwardChatId, enabled := bot.ForwardModeService.GetForwardChat(chatId); enabled {
			if ctx.Message().Text != "" {
				bot.MustSend(forwardChatId,
					fmt.Sprintf("%v %v: %v", ctx.Sender().FirstName, bot.Replies.Says, ctx.Message().Text))
			}
			return nil
		}

		return ctx.Send(bot.Replies.CannotChatYet)
	}
}
