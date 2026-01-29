package bot

import (
	c "memetgbot/internal/core/config"
	l "memetgbot/internal/core/logger"
	"memetgbot/internal/handlers/commands"
	"memetgbot/internal/handlers/message"
	"time"

	"gopkg.in/telebot.v4"
)

func InitBot() {
	var config = c.Config
	var logger = l.Logger

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

	commands.InitCommands(bot)
	message.InitMessagesHandler(bot)

	logger.Info("Bot successfully started!")
	bot.Start()
}
