package main

import (
	b "memetgbot/internal"
	cfg "memetgbot/internal/core/config"
	"memetgbot/internal/core/constants"
	l "memetgbot/internal/core/logger"
	d "memetgbot/internal/db"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/handlers/commands"
	"memetgbot/internal/handlers/message"
	"memetgbot/internal/repo"
	"memetgbot/internal/session"
	"memetgbot/internal/text"
	"memetgbot/pkg/video"
)

func main() {
	config := cfg.MustConfig()
	logger := l.MustLogger(config.IsDebug, l.MustLoggerBot(config.LoggerBotToken), config.AdminID)
	db := d.MustDB(config, logger)
	fsm := fsmManager.New(logger)
	sessionStore := session.NewStore(logger)
	chatRepo := repo.NewChatRepo(db)
	videoService := video.MustNewVideoService(constants.VideoDownloadDirPath, config.YtdlpPath, config.CookiesPath, config.FfmpegPath, logger)
	replies := text.NewReplies()
	bot := b.MustBot(config, fsm, sessionStore, chatRepo, videoService, replies, logger)

	commands.MustInitCommands(bot)
	message.MustInitMessagesHandler(bot)

	logger.Info("Bot successfully started!")
	bot.Start()
}
