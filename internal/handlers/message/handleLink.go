package message

import (
	"context"
	b "memetgbot/internal"
	"memetgbot/internal/core/logger"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/text"
	"memetgbot/pkg/video"
	"strings"

	"gopkg.in/telebot.v4"
)

func handleLink(ctx telebot.Context) error {
	chatId := ctx.Chat().ID

	fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.ProcessingLinkEvent)

	path, name, err := video.VideoService.DownloadVideo(context.Background(), ctx.Message().Text)
	if err != nil {
		if strings.Contains(err.Error(), "URL") {
			logger.Logger.Error(err.Error())
			b.SendWithHandlingErr(chatId, text.Replies.IsNotValidURLError)
			fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			return nil
		}
		logger.Logger.Error(err.Error())
		b.SendWithHandlingErr(chatId, text.Replies.Error)
		fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
		return nil
	}
	a := &telebot.Video{File: telebot.FromDisk(path), FileName: name, CaptionAbove: true, Caption: name}
	b.SendWithHandlingErr(chatId, a)

	fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
	return nil
}
