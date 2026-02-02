package bot

import (
	"fmt"
	"log"
	"memetgbot/internal/core/config"
	"memetgbot/internal/core/logger"
	"memetgbot/internal/feat/forward"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/repo"
	"memetgbot/internal/session"
	"memetgbot/internal/text"
	"memetgbot/models"
	"memetgbot/pkg/video"
	"sync"
	"time"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/react"
)

type Bot struct {
	*telebot.Bot
	Fsm                *fsmManager.FSMState
	SessionStore       *session.Store
	chatRepo           *repo.ChatRepo
	forwardModeRepo    *repo.ForwardMode
	VideoService       *video.VideoService
	ForwardModeService *forward.Service
	Config             *config.AppConfig
	Replies            *text.Replies
	Logger             logger.AppLogger

	chatCache map[int64]*models.Chat
	cacheLock sync.RWMutex
}

func (bot *Bot) MustSend(chatId int64, what interface{}, opts ...interface{}) *telebot.Message {
	msg, err := bot.Send(&telebot.User{ID: chatId}, what, opts...)
	if err != nil {
		bot.Logger.Error(fmt.Sprintf("Error sending message to %v: %v", chatId, err.Error()))
	}
	return msg
}

func (bot *Bot) MustDelete(msg *telebot.Message) {
	err := bot.Delete(msg)
	if err != nil {
		bot.Logger.Error(fmt.Sprintf("Error deleting message to %v: %v", msg.Chat.ID, err.Error()))
	}
}

func (bot *Bot) MustReact(msg *telebot.Message, reaction telebot.Reaction) {
	err := bot.React(&telebot.User{ID: msg.Chat.ID}, msg, react.React(reaction))
	if err != nil {
		bot.Logger.Error(fmt.Sprintf("Error reacting message to %v: %v", msg.Chat.ID, err.Error()))
	}
}

// TODO
func MustEdit() {

}

func (bot *Bot) GetChatCached(chatId int64) (*models.Chat, error) {
	bot.cacheLock.RLock()
	chat, ok := bot.chatCache[chatId]
	bot.cacheLock.RUnlock()

	if ok && chat != nil {
		return chat, nil
	}

	chatDB, err := bot.chatRepo.Get(chatId)
	if err != nil {
		return nil, err
	}

	bot.setChatCache(chatId, &chatDB)

	return &chatDB, nil
}

func (bot *Bot) SaveChat(chat *models.Chat) error {
	if err := bot.chatRepo.Upsert(chat); err != nil {
		return err
	}

	bot.setChatCache(chat.TelegramID, chat)

	return nil
}

func (bot *Bot) setChatCache(chatId int64, chat *models.Chat) {
	bot.cacheLock.Lock()
	defer bot.cacheLock.Unlock()
	bot.chatCache[chatId] = chat
}

func MustBot(
	config *config.AppConfig,
	fsm *fsmManager.FSMState,
	sessionStore *session.Store,
	chatRepo *repo.ChatRepo,
	forwardModeRepo *repo.ForwardMode,
	videoService *video.VideoService,
	forwardModeService *forward.Service,
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

	return &Bot{
		bot,
		fsm,
		sessionStore,
		chatRepo,
		forwardModeRepo,
		videoService,
		forwardModeService,
		config,
		replies,
		logger,
		make(map[int64]*models.Chat),
		sync.RWMutex{},
	}
}
