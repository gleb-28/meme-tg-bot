package main

import (
	b "memetgbot/internal"
	cfg "memetgbot/internal/core/config"
	l "memetgbot/internal/core/logger"
	d "memetgbot/internal/db"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/handlers/commands"
	"memetgbot/internal/handlers/message"
	"memetgbot/internal/repo"
	"memetgbot/internal/text"
	"memetgbot/pkg/video"
)

func main() {
	config := cfg.MustConfig()
	logger := l.MustLogger(config.IsDebug, l.MustLoggerBot(config.LoggerBotToken), config.AdminID)
	db := d.MustDB(config, logger)
	fsm := fsmManager.New(logger)
	chatRepo := repo.NewChatRepo(db)
	videoService := video.MustNewVideoService("./output", "/usr/bin/yt-dlp", logger)
	replies := text.NewReplies()
	bot := b.MustBot(config, fsm, chatRepo, videoService, replies, logger)

	commands.MustInitCommands(bot)
	message.MustInitMessagesHandler(bot)

	logger.Info("Bot successfully started!")
	bot.Start()
}
