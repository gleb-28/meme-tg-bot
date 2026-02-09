package bot

import (
	"fmt"
	"log"
	"memetgbot/internal/core/config"
	"memetgbot/internal/core/constants"
	"memetgbot/internal/core/logger"
	"memetgbot/internal/feat/forward"
	fsmManager "memetgbot/internal/fsm"
	"memetgbot/internal/repo"
	"memetgbot/internal/session"
	"memetgbot/internal/text"
	"memetgbot/model"
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

	chatCache map[int64]*model.Chat
	cacheLock sync.RWMutex
}

func (bot *Bot) MustSend(chatId int64, what interface{}, opts ...interface{}) *telebot.Message {
	msg, err := bot.Send(&telebot.User{ID: chatId}, what, opts...)
	if err != nil {
		bot.Logger.Error(fmt.Sprintf("Error sending message to %v: %v", chatId, err.Error()))
	}
	return msg
}

func (bot *Bot) MustSendAlbum(chatID int64, album telebot.Album) []telebot.Message {
	msg, err := bot.Bot.SendAlbum(&telebot.Chat{ID: chatID}, album)
	if err != nil {
		bot.Logger.Error(fmt.Sprintf("Error sending album to %v: %v", chatID, err.Error()))
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
func MustEdit() {}

func (bot *Bot) ForwardAnyMessage(ctx telebot.Context, forwardChatId int64) {
	msg := ctx.Message()
	sender := ctx.Sender().FirstName

	saysPrefix := fmt.Sprintf("%s %s", sender, bot.Replies.Says)
	sendsPrefix := fmt.Sprintf("%s %s", sender, bot.Replies.Sends)

	if msg.AlbumID != "" {
		bot.handleAlbumGroup(ctx, forwardChatId)
		return
	}

	// ===== TEXT =====
	if msg.Text != "" {
		bot.MustSend(
			forwardChatId,
			fmt.Sprintf("%s: %s", saysPrefix, msg.Text),
		)
		bot.MustReact(msg, react.ThumbUp)
		return
	}

	// ===== PHOTO =====
	if msg.Photo != nil {
		caption := sendsPrefix
		if msg.Caption != "" {
			caption = fmt.Sprintf("%s: %s", sendsPrefix, msg.Caption)
		}

		photo := &telebot.Photo{
			File:    msg.Photo.File,
			Caption: caption,
		}

		bot.MustSend(forwardChatId, photo)
		bot.MustReact(msg, react.ThumbUp)
		return
	}

	// ===== VIDEO =====
	if msg.Video != nil {
		caption := sendsPrefix
		if msg.Caption != "" {
			caption = fmt.Sprintf("%s: %s", sendsPrefix, msg.Caption)
		}

		vid := &telebot.Video{
			File:    msg.Video.File,
			Caption: caption,
		}

		bot.MustSend(forwardChatId, vid)
		bot.MustReact(msg, react.ThumbUp)
		return
	}

	// ===== AUDIO =====
	if msg.Audio != nil {
		audio := &telebot.Audio{
			File:    msg.Audio.File,
			Caption: sendsPrefix,
		}
		bot.MustSend(forwardChatId, audio)
		bot.MustReact(msg, react.ThumbUp)
		return
	}

	// ===== DOCUMENT =====
	if msg.Document != nil {
		doc := &telebot.Document{
			File:    msg.Document.File,
			Caption: sendsPrefix,
		}
		bot.MustSend(forwardChatId, doc)
		bot.MustReact(msg, react.ThumbUp)
		return
	}

	// ===== STICKER =====
	if msg.Sticker != nil {
		bot.MustSend(forwardChatId, sendsPrefix+":")
		bot.MustSend(forwardChatId, msg.Sticker)
		bot.MustReact(msg, react.ThumbUp)
		return
	}

	// ===== VOICE =====
	if msg.Voice != nil {
		voice := &telebot.Voice{
			File: msg.Voice.File,
		}
		bot.MustSend(forwardChatId, saysPrefix+":")
		bot.MustSend(forwardChatId, voice)
		bot.MustReact(msg, react.ThumbUp)
		return
	}
}

func (bot *Bot) handleAlbumGroup(ctx telebot.Context, forwardChatId int64) {
	msg := ctx.Message()
	userID := ctx.Chat().ID

	batch, exists := bot.SessionStore.GetMediaBatch(userID)
	albumTimeout := 600 * time.Millisecond

	if !exists {
		batch = &session.MediaBatch{}

		batch.Timer = time.AfterFunc(albumTimeout, func() {
			bot.SessionStore.DeleteMediaBatch(userID)

			if len(batch.Items) == 0 {
				return
			}

			album := make(telebot.Album, len(batch.Items))
			for i, item := range batch.Items {
				album[i] = item
			}

			if batch.Caption != "" {
				if len(batch.Caption) < constants.TelegramMaxCaptionLen {
					album.SetCaption(batch.Caption)
				} else {
					defer bot.MustSend(forwardChatId, batch.Caption)
				}
			}

			bot.MustSendAlbum(forwardChatId, album)
			bot.MustReact(msg, react.ThumbUp)
		})

		bot.SessionStore.SetMediaBatch(userID, batch)
	} else {
		batch.Timer.Reset(albumTimeout)
	}

	// ===== ADD MEDIA =====

	var item telebot.Inputtable

	if msg.Photo != nil {
		item = &telebot.Photo{
			File: msg.Photo.File,
		}
	}

	if msg.Video != nil {
		item = &telebot.Video{
			File: msg.Video.File,
		}
	}

	if msg.Audio != nil {
		item = &telebot.Audio{
			File: msg.Audio.File,
		}
	}

	if msg.Document != nil {
		item = &telebot.Document{
			File: msg.Document.File,
		}
	}

	if item != nil {
		batch.Items = append(batch.Items, item)
	}

	// ===== CAPTION ONCE =====

	if msg.Caption != "" && batch.Caption == "" {
		batch.Caption = fmt.Sprintf(
			"%s %s: %s",
			ctx.Sender().FirstName,
			bot.Replies.Sends,
			msg.Caption,
		)
	}
}

func (bot *Bot) GetChatCached(chatId int64) (*model.Chat, error) {
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

func (bot *Bot) SaveChat(chat *model.Chat) error {
	if err := bot.chatRepo.Upsert(chat); err != nil {
		return err
	}

	bot.setChatCache(chat.TelegramID, chat)

	return nil
}

func (bot *Bot) setChatCache(chatId int64, chat *model.Chat) {
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
		make(map[int64]*model.Chat),
		sync.RWMutex{},
	}
}
