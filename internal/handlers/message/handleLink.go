package message

import (
	"context"
	b "memetgbot/internal"
	"memetgbot/internal/session"
	"memetgbot/pkg/utils"
	"strconv"
	"strings"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/react"
)

func createHandleLink(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		userMsg := ctx.Message()
		processingLinkKey := strconv.FormatInt(int64(userMsg.ID), 10)

		botMsg := bot.MustSend(chatId, bot.Replies.Downloading)

		bot.SessionStore.AddProcessingLink(chatId, processingLinkKey, &session.ProcessingLink{
			UserMsg: userMsg,
			BotMsg:  botMsg,
		})

		path, name, err := bot.VideoService.DownloadVideo(context.Background(), ctx.Message().Text)
		if err != nil {
			if strings.Contains(err.Error(), "URL") {
				bot.Logger.Error(err.Error())
				bot.MustSend(chatId, bot.Replies.IsNotValidURLError)
			} else {
				bot.Logger.Error(err.Error())
				bot.MustSend(chatId, bot.Replies.Error)
			}

			bot.MustDelete(botMsg)
			bot.SessionStore.RemoveProcessingLink(chatId, processingLinkKey)

			return nil
		}

		cleanFileName := utils.RemoveSaltFromFileName(utils.RemoveCompressedSuffix(name))
		video := &telebot.Video{File: telebot.FromDisk(path), FileName: cleanFileName, CaptionAbove: true, Caption: cleanFileName}
		bot.MustSend(chatId, video)

		bot.MustReact(userMsg, react.ThumbUp)
		bot.MustDelete(botMsg)

		err = bot.VideoService.DeleteVideoByName(name)

		bot.SessionStore.RemoveProcessingLink(chatId, processingLinkKey)

		return nil
	}
}
