package session

import (
	"fmt"
	"memetgbot/internal/core/logger"
	"reflect"
	"sync"

	"gopkg.in/telebot.v4"
)

type Session struct {
	LastLinkProcessingMessage *telebot.Message // user msg with link
	LastBotDownloadingMessage *telebot.Message // bot reply with "downloading" state
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

	session = &Session{}
	store.sessions[chatID] = session
	return session
}

func (store *Store) SetField(chatID int64, fieldName string, value interface{}) {
	session := store.Get(chatID)

	store.mu.Lock()
	defer store.mu.Unlock()

	v := reflect.ValueOf(session).Elem()
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		store.logger.Error(fmt.Sprintf("field %s not found in Session", fieldName))
		return
	}
	if !f.CanSet() {
		store.logger.Error(fmt.Sprintf("field %s cannot be set", fieldName))
		return
	}

	val := reflect.ValueOf(value)

	if val.Type().ConvertibleTo(f.Type()) {
		f.Set(val.Convert(f.Type()))
	}
	store.logger.Error(fmt.Sprintf("cannot set field %s with value of type %s", fieldName, val.Type()))
}
