package forward

import (
	"fmt"
	"memetgbot/internal/core/logger"
	"memetgbot/internal/repo"
	"memetgbot/internal/session"
)

type Service struct {
	forwardModeRepo *repo.ForwardMode
	session         *session.Store
	logger          logger.AppLogger
}

func NewForwardModeService(forwardModeRepo *repo.ForwardMode, sessionStore *session.Store, logger logger.AppLogger) *Service {
	return &Service{
		forwardModeRepo: forwardModeRepo,
		session:         sessionStore,
		logger:          logger,
	}
}

func (s *Service) EnableForwardMode(userID int64, forwardChatId int64) error {
	s.ensureForwardModeLoaded(userID)
	err := s.forwardModeRepo.Enable(userID, forwardChatId)
	if err != nil {
		return err
	}

	s.session.EnableForwardMode(userID, forwardChatId)

	return nil
}

func (s *Service) DisableForwardMode(userID int64) error {
	s.ensureForwardModeLoaded(userID)
	err := s.forwardModeRepo.Disable(userID)
	if err != nil {
		return err
	}
	s.session.DisableForwardMode(userID)

	return nil
}

func (s *Service) GetForwardChat(userID int64) (chatID int64, enabled bool) {
	s.ensureForwardModeLoaded(userID)
	enabled, chatID = s.session.GetForwardMode(userID)
	return
}

func (s *Service) IsForwardEnabled(userID int64) bool {
	s.ensureForwardModeLoaded(userID)
	enabled, _ := s.session.GetForwardMode(userID)
	return enabled
}

func (s *Service) ensureForwardModeLoaded(userID int64) {
	if s.session.IsForwardModeLoaded(userID) {
		return
	}

	fm, ok, err := s.forwardModeRepo.TryGet(userID)
	if err == nil && ok {
		s.session.SetForwardMode(userID, fm.IsEnabled, fm.ChatID)
		return
	}

	if err != nil {
		s.logger.Error(fmt.Sprintf("Error while ensureForwardModeLoaded %v: %v", userID, err.Error()))
	}
}
