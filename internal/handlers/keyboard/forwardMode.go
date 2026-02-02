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

		if bot.ForwardModeService.IsForwardEnabled(chatId) {
			bot.MustSend(chatId, bot.Replies.ForwardingModeIsAlreadyEnabled)
			return nil
		}

		bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.AwaitingForwardChatEvent)

		bot.MustSend(chatId, bot.Replies.ForwardMsgFromGroup)

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

		return ctx.Edit(bot.Replies.ForwardingModeIsDisabled, ForwardModeKeyboard(false, bot.Replies))
	}
}

func ForwardModeKeyboard(isEnabled bool, replies *text.Replies) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{IsPersistent: true, OneTimeKeyboard: true}

	var toggleBtn telebot.Btn
	if isEnabled {
		toggleBtn = kb.Data(replies.DisableForwardingMode, constants.DisableForwardModeButtonCallbackQuery)
	} else {
		toggleBtn = kb.Data(replies.EnableForwardingMode, constants.EnableForwardModeButtonCallbackQuery)
	}

	kb.Inline(kb.Row(toggleBtn))

	return kb
}
