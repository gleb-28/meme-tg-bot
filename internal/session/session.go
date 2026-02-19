package session

import (
	"context"
	"memetgbot/internal/core/logger"
	"sync"
	"time"

	"gopkg.in/telebot.v4"
)

type Session struct {
	ProcessingLinks      map[string]*ProcessingLink
	MediaBatches         map[string]*MediaBatch
	ForwardModeLoaded    bool
	ForwardModeIsEnabled bool
	ForwardChatID        int64
	Authorized           bool
	ExpiresAt            time.Time
}

type ProcessingLink struct {
	UserMsg *telebot.Message // user msg with link
	BotMsg  *telebot.Message // bot reply with "downloading" state
}

type MediaItem struct {
	MessageID int
	Input     telebot.Inputtable
}

type MediaBatch struct {
	Items   []MediaItem
	Caption string
	Timer   *time.Timer
}

type Store struct {
	mu       sync.RWMutex
	sessions map[int64]*Session
	ttl      time.Duration
	logger   logger.AppLogger
}

func NewStore(ttl time.Duration, logger logger.AppLogger) *Store {
	return &Store{
		sessions: make(map[int64]*Session),
		ttl:      ttl,
		logger:   logger,
	}
}

func (store *Store) Get(chatID int64) *Session {
	store.mu.Lock()
	defer store.mu.Unlock()

	session, ok := store.sessions[chatID]
	if !ok {
		session = &Session{
			ProcessingLinks: make(map[string]*ProcessingLink),
			MediaBatches:    make(map[string]*MediaBatch),
		}
		store.sessions[chatID] = session
	}

	store.ensureExpiryLocked(session)

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
	session.ForwardChatID = forwardChatID
}

func (store *Store) DisableForwardMode(chatID int64) {
	session := store.Get(chatID)
	store.mu.Lock()
	defer store.mu.Unlock()
	session.ForwardModeIsEnabled = false
}

func (store *Store) GetForwardMode(chatID int64) (isEnabled bool, forwardChatID int64) {
	session := store.Get(chatID)
	store.mu.RLock()
	defer store.mu.RUnlock()
	return session.ForwardModeIsEnabled, session.ForwardChatID
}

func (store *Store) SetForwardMode(chatID int64, isEnabled bool, forwardChatID int64) {
	session := store.Get(chatID)
	store.mu.Lock()
	defer store.mu.Unlock()
	session.ForwardModeIsEnabled = isEnabled
	session.ForwardChatID = forwardChatID
	session.ForwardModeLoaded = true
}

// --- Media batches (albums) ---

// GetOrCreateMediaBatch returns existing batch for albumID or creates it atomically.
// If created, init is called before releasing the lock so the batch is fully
// initialized before use by other goroutines.
func (store *Store) GetOrCreateMediaBatch(chatID int64, albumID string, init func(*MediaBatch)) (*MediaBatch, bool) {
	session := store.Get(chatID)
	store.mu.Lock()
	defer store.mu.Unlock()
	batch, ok := session.MediaBatches[albumID]
	if !ok {
		batch = &MediaBatch{}
		if init != nil {
			init(batch)
		}
		session.MediaBatches[albumID] = batch
	}
	return batch, ok
}

// DeleteMediaBatch removes and returns whether it existed.
func (store *Store) DeleteMediaBatch(chatID int64, albumID string) bool {
	session := store.Get(chatID)
	store.mu.Lock()
	defer store.mu.Unlock()
	if _, ok := session.MediaBatches[albumID]; ok {
		delete(session.MediaBatches, albumID)
		return true
	}
	return false
}

// Delete removes session and stops timers to avoid leaks.
func (store *Store) Delete(chatID int64) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.deleteLocked(chatID)
}

// MarkAuthorized stops expiration for an authenticated chat.
func (store *Store) MarkAuthorized(chatID int64) {
	store.mu.Lock()
	defer store.mu.Unlock()

	session, ok := store.sessions[chatID]
	if !ok {
		return
	}

	session.Authorized = true
	session.ExpiresAt = time.Time{}
}

func (store *Store) ensureExpiryLocked(session *Session) {
	if session.Authorized || store.ttl <= 0 {
		session.ExpiresAt = time.Time{}
		return
	}

	session.ExpiresAt = time.Now().Add(store.ttl)
}

func (store *Store) StartCleanupWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				store.cleanupExpired()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (store *Store) cleanupExpired() {
	now := time.Now()

	store.mu.Lock()
	defer store.mu.Unlock()

	for chatID, session := range store.sessions {
		if session.Authorized {
			continue
		}

		if !session.ExpiresAt.IsZero() && now.After(session.ExpiresAt) {
			store.deleteLocked(chatID)
		}
	}
}

func (store *Store) deleteLocked(chatID int64) {
	session, ok := store.sessions[chatID]
	if !ok {
		return
	}

	for _, batch := range session.MediaBatches {
		if batch.Timer != nil {
			batch.Timer.Stop()
		}
	}

	delete(store.sessions, chatID)
}
