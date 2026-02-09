package keyboard

import (
	"context"
	"fmt"
	b "memetgbot/internal"
	"memetgbot/internal/core/constants"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/text"

	"gopkg.in/telebot.v4"
)

func createForwardModeEnabledHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Sender().ID
		forwardChatId, forwardModeIsEnabled := bot.ForwardModeService.GetForwardChat(chatId)

		if forwardModeIsEnabled {
			bot.MustSend(chatId, bot.Replies.ForwardingModeIsAlreadyEnabled)
			return nil
		}

		if forwardChatId == 0 {
			enableAwaitingForwardChatEvent(bot, chatId)
			return nil
		}

		forwardChat, err := bot.ChatByID(forwardChatId)
		if err != nil {
			enableAwaitingForwardChatEvent(bot, chatId)
			return nil
		}

		member, err := bot.ChatMemberOf(forwardChat, bot.Me)
		if err != nil || (member.Role != telebot.Administrator && member.Role != telebot.Creator) {
			enableAwaitingForwardChatEvent(bot, chatId)
			return nil
		}

		bot.MustEdit(ctx.Message(),
			fmt.Sprintf("%v: \"%v\"\n%v",
				bot.Replies.YouHaveThisForwardChat, forwardChat.Title, bot.Replies.UseItAgain),
			usePreviousForwardChatKeyboard(bot.Replies))

		return nil
	}
}

func enableAwaitingForwardChatEvent(bot *b.Bot, chatId int64) {
	bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.AwaitingForwardChatEvent)

	bot.MustSend(chatId, bot.Replies.ForwardMsgFromGroup)
}

func createChangeForwardChatHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Sender().ID

		bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.AwaitingForwardChatEvent)

		bot.MustSend(chatId, bot.Replies.ForwardMsgFromGroup)

		return nil
	}
}

func createUsePrevForwardChatHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Sender().ID

		forwardChatId, _ := bot.ForwardModeService.GetForwardChat(chatId)
		forwardChat, err := bot.ChatByID(forwardChatId)
		if err != nil {
			bot.Logger.Error(fmt.Sprintf("Error finding chat by id %v: %v", forwardChatId, err.Error()))
			bot.MustSend(chatId, bot.Replies.Error)
			return nil
		}

		err = bot.ForwardModeService.EnableForwardMode(chatId, forwardChatId)
		if err != nil {
			bot.Logger.Error(fmt.Sprintf("Error enabling forward mode to %v: %v", chatId, err.Error()))
			bot.MustSend(chatId, bot.Replies.Error)
			return nil
		}

		bot.MustSend(chatId, fmt.Sprintf("%v\"%v\"", bot.Replies.SuccessEnablingForwardMode, forwardChat.Title))
		bot.MustEdit(ctx.Message(), bot.Replies.ForwardingModeIsEnabled, ForwardModeKeyboard(true, bot.Replies))

		return nil
	}
}

func createForwardModeDisabledHandler(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Sender().ID

		if !bot.ForwardModeService.IsForwardEnabled(chatId) {
			bot.MustSend(chatId, bot.Replies.ForwardingModeIsAlreadyDisabled)
			return nil
		}

		err := bot.ForwardModeService.DisableForwardMode(chatId)
		if err != nil {
			bot.MustSend(chatId, bot.Replies.Error)
			bot.Logger.Error(fmt.Sprintf("Error disabling forward mode to %v: %v", chatId, err.Error()))
		}

		bot.MustEdit(ctx.Message(), bot.Replies.ForwardingModeIsDisabled, ForwardModeKeyboard(false, bot.Replies))
		return nil
	}
}

func ForwardModeKeyboard(isEnabled bool, replies *text.Replies) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{ResizeKeyboard: true}

	var toggleBtn telebot.Btn
	if isEnabled {
		toggleBtn = kb.Data(replies.DisableForwardingMode, constants.DisableForwardModeButtonCallbackQuery)
	} else {
		toggleBtn = kb.Data(replies.EnableForwardingMode, constants.EnableForwardModeButtonCallbackQuery)
	}

	kb.Inline(kb.Row(toggleBtn))

	return kb
}

func usePreviousForwardChatKeyboard(replies *text.Replies) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{ResizeKeyboard: true}

	btnUse := kb.Data(replies.UsePrevForwardChat, constants.UsePrevForwardChatButtonCallbackQuery)
	btnChange := kb.Data(replies.ChangePrevForwardChat, constants.ChangeForwardChatButtonCallbackQuery)

	kb.Inline(
		kb.Row(btnUse),
		kb.Row(btnChange),
	)

	return kb
}
