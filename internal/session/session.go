package session

import (
	"memetgbot/internal/core/logger"
	"sync"

	"gopkg.in/telebot.v4"
)

type Session struct {
	ProcessingLinks      map[string]*ProcessingLink
	ForwardModeLoaded    bool
	ForwardModeIsEnabled bool
	ForwardChatId        int64
}

type ProcessingLink struct {
	UserMsg *telebot.Message // user msg with link
	BotMsg  *telebot.Message // bot reply with "downloading" state
}

type Store struct {
	mu       sync.RWMutex
	sessions map[int64]*Session
	logger   logger.AppLogger
}

func NewStore(logger logger.AppLogger) *Store {
	return &Store{
		sessions: make(map[int64]*Session),
		logger:   logger,
	}
}

func (store *Store) Get(chatID int64) *Session {
	store.mu.RLock()
	session, ok := store.sessions[chatID]
	store.mu.RUnlock()

	if ok {
		return session
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	session = &Session{
		ProcessingLinks: make(map[string]*ProcessingLink),
	}
	store.sessions[chatID] = session
	return session
}

// --- ProcessingLinks ---

func (store *Store) AddProcessingLink(chatID int64, key string, link *ProcessingLink) {
	session := store.Get(chatID)

	store.mu.Lock()
	defer store.mu.Unlock()

	session.ProcessingLinks[key] = link
}

func (store *Store) GetProcessingLink(chatID int64, key string) (*ProcessingLink, bool) {
	session := store.Get(chatID)

	store.mu.RLock()
	defer store.mu.RUnlock()

	pl, ok := session.ProcessingLinks[key]
	return pl, ok
}

func (store *Store) RemoveProcessingLink(chatID int64, key string) {
	session := store.Get(chatID)

	store.mu.Lock()
	defer store.mu.Unlock()

	delete(session.ProcessingLinks, key)
}

// --- Forward Mode ---

func (store *Store) IsForwardModeLoaded(chatID int64) bool {
	session := store.Get(chatID)

	store.mu.RLock()
	defer store.mu.RUnlock()

	return session.ForwardModeLoaded
}

func (store *Store) EnableForwardMode(chatID int64, forwardChatID int64) {
	session := store.Get(chatID)
	store.mu.Lock()
	defer store.mu.Unlock()
	session.ForwardModeIsEnabled = true
	session.ForwardChatId = forwardChatID
	session.ForwardModeLoaded = true
}

func (store *Store) DisableForwardMode(chatID int64) {
	session := store.Get(chatID)
	store.mu.Lock()
	defer store.mu.Unlock()
	session.ForwardModeIsEnabled = false
	session.ForwardChatId = 0
	session.ForwardModeLoaded = true
}

func (store *Store) GetForwardMode(chatID int64) (isEnabled bool, forwardChatID int64) {
	session := store.Get(chatID)
	store.mu.RLock()
	defer store.mu.RUnlock()
	return session.ForwardModeIsEnabled, session.ForwardChatId
}
