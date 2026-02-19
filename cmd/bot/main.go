package main

import (
	"context"
	b "memetgbot/internal"
	cfg "memetgbot/internal/core/config"
	"memetgbot/internal/core/constants"
	l "memetgbot/internal/core/logger"
	d "memetgbot/internal/db"
	"memetgbot/internal/feat/forward"
	"memetgbot/internal/feat/media"
	"memetgbot/internal/feat/media/instagram"
	"memetgbot/internal/feat/media/video"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/handler/commands"
	"memetgbot/internal/handler/keyboard"
	"memetgbot/internal/handler/message"
	"memetgbot/internal/repo"
	"memetgbot/internal/session"
	"memetgbot/internal/text"
)

func main() {
	config := cfg.MustConfig()
	logger := l.MustLogger(config.IsDebug, l.MustLoggerBot(config.LoggerBotToken), config.AdminID)

	db := d.MustDB(config, logger)

	fsm := fsmManager.New(logger)
	sessionStore := session.NewStore(constants.NonAuthSessionTTL, logger)
	sessionStore.StartCleanupWorker(context.Background(), constants.NonAuthSessionTTL/2)

	chatRepo := repo.NewChatRepo(db)
	forwardModeRepo := repo.NewForwardModeRepo(db)

	videoService := video.MustNewVideoService(constants.VideoDownloadDirPath, config.YtdlpPath, config.CookiesPath, config.FfmpegPath, logger)
	instaService := instagram.NewService(videoService, instagram.NewImageService(constants.VideoDownloadDirPath, config.CookiesPath, logger))
	mediaService := media.NewService(instaService, videoService)
	forwardModeService := forward.NewForwardModeService(forwardModeRepo, sessionStore, logger)

	replies := text.NewReplies()
	bot := b.MustBot(config, fsm, sessionStore, chatRepo, forwardModeRepo, mediaService, forwardModeService, replies, logger)

	commands.MustInitCommandsHandler(bot)
	message.MustInitMessagesHandler(bot)
	keyboard.MustInitKeyboardHandler(bot)

	logger.Info("Bot successfully started!")
	bot.Start()
}
