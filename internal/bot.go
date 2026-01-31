package bot

import (
	"fmt"
	"log"
	"memetgbot/internal/core/config"
	"memetgbot/internal/core/logger"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/repo"
	"memetgbot/internal/text"
	"memetgbot/pkg/video"
	"time"

	"gopkg.in/telebot.v4"
)

type Bot struct {
	*telebot.Bot
	Fsm          *fsmManager.FSMState
	ChatRepo     *repo.ChatRepo
	VideoService *video.VideoService
	Config       *config.AppConfig
	Replies      *text.Replies
	Logger       logger.AppLogger
}

func (bot Bot) MustSend(chatId int64, what interface{}, opts ...interface{}) {
	_, err := bot.Send(&telebot.User{ID: chatId}, what, opts...)
	if err != nil {
		bot.Logger.Error(fmt.Sprintf("Error sending message to %v: %v", chatId, err.Error()))
	}
}

func MustBot(
	config *config.AppConfig,
	fsm *fsmManager.FSMState,
	chatRepo *repo.ChatRepo,
	videoService *video.VideoService,
	replies *text.Replies,
	logger logger.AppLogger) *Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:     config.TgBotToken,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeHTML,
		OnError: func(err error, context telebot.Context) {
			logger.Error(err.Error())
		},
	})
	if err != nil {
		log.Fatal("Error creating bot:", err.Error())
	}

	return &Bot{bot, fsm, chatRepo, videoService, config, replies, logger}
}
