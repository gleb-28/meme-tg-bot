package message

import (
	"context"
	"fmt"
	b "memetgbot/internal"
	"memetgbot/internal/session"
	"memetgbot/model"
	"memetgbot/pkg/utils"
	"strconv"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/react"
)

func createHandleLink(bot *b.Bot) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		chatId := ctx.Chat().ID
		userMsg := ctx.Message()
		userFirstName := ctx.Sender().FirstName
		processingLinkKey := strconv.FormatInt(int64(userMsg.ID), 10)

		botMsg := bot.MustSend(chatId, bot.Replies.Downloading)

		bot.SessionStore.AddProcessingLink(chatId, processingLinkKey, &session.ProcessingLink{
			UserMsg: userMsg,
			BotMsg:  botMsg,
		})

		reqCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		res, err := bot.MediaService.Extract(reqCtx, ctx.Message().Text)
		if err != nil {
			if strings.Contains(err.Error(), "URL") {
				bot.Logger.Error(err.Error())
				bot.MustSend(chatId, bot.Replies.IsNotValidURLError)
			} else {
				bot.Logger.Error(err.Error())
				bot.MustSend(chatId, bot.Replies.Error)
			}
			return nil
		}

		if forwardChatId, enabled := bot.ForwardModeService.GetForwardChat(chatId); enabled {
			sendMedia(bot, forwardChatId, res, fmt.Sprintf(bot.Replies.ValueFromName, utils.RemoveCompressedSuffix(res.Files[0].Name), userFirstName))
		} else {
			sendMedia(bot, chatId, res, utils.RemoveCompressedSuffix(res.Files[0].Name))
		}

		bot.MustReact(userMsg, react.ThumbUp)
		bot.MustDelete(botMsg)

		bot.SessionStore.RemoveProcessingLink(chatId, processingLinkKey)

		return nil
	}
}

func sendMedia(bot *b.Bot, chatId int64, media *model.MediaResult, caption string) {
	switch media.Type {
	case model.MediaAlbum:
		bot.MustSendAlbum(chatId, album(media, caption))
	case model.MediaVideo:
		bot.MustSend(chatId, video(media.Files[0].Path, media.Files[0].Name, caption))
	}
}

func video(path string, name string, caption string) *telebot.Video {
	return &telebot.Video{
		File:     telebot.FromDisk(path),
		FileName: name, CaptionAbove: true,
		Caption: caption,
	}
}

func album(media *model.MediaResult, caption string) telebot.Album {
	album := make(telebot.Album, len(media.Files))
	for i, f := range media.Files {
		album[i] = &telebot.Photo{
			File: telebot.FromDisk(f.Path),
		}
	}
	album.SetCaption(caption)
	return album
}
