package message

import (
	"context"
	"fmt"
	b "memetgbot/internal"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/handler/keyboard"

	"gopkg.in/telebot.v4"
)

func createValidateForwardChatHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Sender().ID

		forwardChatOrigin := ctx.Message().Origin
		if forwardChatOrigin == nil {
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			bot.MustSend(chatId, bot.Replies.UnableToFindTheGroup)
			return nil
		}
		forwardChat := forwardChatOrigin.SenderChat
		if forwardChat == nil {
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			bot.MustSend(chatId, bot.Replies.UnableToFindTheGroup)
			return nil
		}

		member, err := bot.ChatMemberOf(forwardChat, bot.Me)
		if err != nil || (member.Role != telebot.Administrator && member.Role != telebot.Creator) {
			bot.MustSend(chatId, bot.Replies.IMustBeAdminInTheGroup)
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			return nil
		}

		err = bot.ForwardModeService.EnableForwardMode(chatId, forwardChat.ID)
		if err != nil {
			bot.Logger.Error(fmt.Sprintf("Error enabling forward mode to %v: %v", chatId, err.Error()))
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			bot.MustSend(chatId, bot.Replies.Error)
			return nil
		}

		bot.MustSend(chatId, fmt.Sprintf("%v\"%v\"", bot.Replies.SuccessEnablingForwardMode, forwardChat.Title))
		bot.MustEdit(ctx.Message(), bot.Replies.ForwardingModeIsEnabled, keyboard.ForwardModeKeyboard(true, bot.Replies))

		bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)

		return nil
	}
}
