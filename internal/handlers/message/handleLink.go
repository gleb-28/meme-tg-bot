package message

import (
	"context"
	"memetgbot/internal/core/logger"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/repo"
	"memetgbot/internal/text"
	"memetgbot/pkg/video"
	"strings"

	"gopkg.in/telebot.v4"
)

func handleLink(ctx telebot.Context) error {
	chatId := ctx.Chat().ID

	chat, err := repo.Chat.Get(chatId)
	if err != nil {
		logger.Logger.Error(err.Error())
		return ctx.Send(text.Replies.Error)
	}
	if chat.TelegramID == 0 {
		return nil
	}

	fsm := fsmManager.FSM.GetFSMForUser(chatId)
	err = fsm.Event(context.Background(), fsmManager.ProcessingLinkEvent)

	path, name, err := video.VideoService.DownloadVideo(ctx.Message().Text)
	if err != nil {
		if strings.Contains(err.Error(), "is not a valid URL") {
			logger.Logger.Error(err.Error())
			ctx.Send(text.Replies.IsNotValidURLError)
			return setInitialState(ctx, chatId)
		}
		logger.Logger.Error(err.Error())
		ctx.Send(text.Replies.Error)
		return setInitialState(ctx, chatId)
	}
	a := &telebot.Video{File: telebot.FromDisk(path), FileName: name, CaptionAbove: true, Caption: name}
	err = ctx.Send(a)
	if err != nil {
		logger.Logger.Error(err.Error())
		ctx.Send(text.Replies.Error)
		return setInitialState(ctx, chatId)
	}

	return setInitialState(ctx, chatId)
}
