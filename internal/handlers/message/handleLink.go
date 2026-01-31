package message

import (
	"context"
	b "memetgbot/internal"
	"memetgbot/internal/core/logger"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/text"
	"memetgbot/pkg/utils"
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
	cleanFileName := utils.RemoveSaltFromFileName(name)
	a := &telebot.Video{File: telebot.FromDisk(path), FileName: cleanFileName, CaptionAbove: true, Caption: cleanFileName}
	b.SendWithHandlingErr(chatId, a)
	err = video.VideoService.DeleteVideoByName(name)
	if err != nil {
		logger.Logger.Error("Error deleting file " + name + err.Error())
	}

	fsmManager.FSM.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
	return nil
}
