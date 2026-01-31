package bot

import (
	"fmt"
	c "memetgbot/internal/core/config"
	l "memetgbot/internal/core/logger"
	"time"

	"gopkg.in/telebot.v4"
)

var logger = l.Logger

func createBot() *telebot.Bot {
	var config = c.Config

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     config.TgBotToken,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeHTML,
		OnError: func(err error, context telebot.Context) {
			logger.Error(err.Error())
		},
	})
	if err != nil {
		panic("Error creating bot: " + err.Error())
	}

	return bot
}

var Bot = createBot()

func SendWithHandlingErr(chatId int64, what interface{}, opts ...interface{}) {
	_, err := Bot.Send(&telebot.User{ID: chatId}, what, opts...)
	if err != nil {
		logger.Error(fmt.Sprintf("Error sending message to %v: %v", chatId, err.Error()))
	}
}
