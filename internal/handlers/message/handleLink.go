package message

import (
	"context"
	b "memetgbot/internal"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/pkg/utils"
	"strings"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/react"
)

func createHandleLink(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID

		botDownloadingMessage := bot.MustSend(chatId, bot.Replies.Downloading)
		session := bot.SessionStore.Get(chatId)
		session.LastLinkProcessingMessage = ctx.Message()
		session.LastBotDownloadingMessage = botDownloadingMessage
		bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.ProcessingLinkEvent)

		path, name, err := bot.VideoService.DownloadVideo(context.Background(), ctx.Message().Text)
		if err != nil {
			if strings.Contains(err.Error(), "URL") {
				bot.Logger.Error(err.Error())
				bot.MustSend(chatId, bot.Replies.IsNotValidURLError)
				bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
				return nil
			}
			bot.Logger.Error(err.Error())
			bot.MustSend(chatId, bot.Replies.Error)
			bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
			return nil
		}
		cleanFileName := utils.RemoveSaltFromFileName(utils.RemoveCompressedSuffix(name))
		a := &telebot.Video{File: telebot.FromDisk(path), FileName: cleanFileName, CaptionAbove: true, Caption: cleanFileName}
		bot.MustSend(chatId, a)

		err = bot.React(
			&telebot.User{ID: chatId},
			session.LastLinkProcessingMessage,
			react.React(react.ThumbUp),
		)
		if err != nil {
			bot.Logger.Error(err.Error())
		}

		err = bot.Delete(session.LastBotDownloadingMessage)
		if err != nil {
			bot.Logger.Error(err.Error())
		}

		err = bot.VideoService.DeleteVideoByName(name)
		if err != nil {
			bot.Logger.Error("Error deleting file " + name + err.Error())
		}

		bot.Fsm.UserEvent(context.Background(), chatId, fsmManager.InitialEvent)
		return nil
	}
}
