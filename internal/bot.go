package bot

import (
	c "memetgbot/internal/core/config"
	l "memetgbot/internal/core/logger"
	"memetgbot/internal/handlers/commands"
	"time"

	t "gopkg.in/telebot.v4"
)

func InitBot() {
	var config = c.Config
	var logger = l.Logger

	bot, err := t.NewBot(t.Settings{
		Token:     config.TgBotToken,
		Poller:    &t.LongPoller{Timeout: 10 * time.Second},
		ParseMode: t.ModeHTML,
		OnError: func(err error, context t.Context) {
			logger.Error(err.Error())
		},
	})
	if err != nil {
		panic("Error creating bot: " + err.Error())
	}

	commands.InitCommands(bot)

	logger.Info("Bot successfully started!")
	bot.Start()
}
