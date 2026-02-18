package repo

import (
	"errors"
	"memetgbot/model"

	"gorm.io/gorm"
)

type ForwardMode struct {
	db *gorm.DB
}

func NewForwardModeRepo(db *gorm.DB) *ForwardMode {
	return &ForwardMode{db: db}
}

func (r *ForwardMode) TryGet(userID int64) (*model.ForwardMode, bool, error) {
	var result model.ForwardMode

	err := r.db.Where("user_id = ?", userID).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &result, true, nil
}

func (r *ForwardMode) Enable(userID int64, chatId int64) error {
	fm := model.ForwardMode{
		UserID:    userID,
		IsEnabled: true,
		ChatID:    chatId,
	}
	return r.db.Model(&model.ForwardMode{}).
		Where("user_id = ?", userID).
		Assign(fm).
		FirstOrCreate(&fm).
		Error
}

func (r *ForwardMode) Disable(userID int64) error {
	return r.db.
		Model(&model.ForwardMode{}).
		Where("user_id = ?", userID).
		Update("is_enabled", false).
		Error
}

func (r *ForwardMode) IsEnabled(userID int64) (bool, error) {
	fm, found, err := r.TryGet(userID)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}
	return fm.IsEnabled, nil
}
