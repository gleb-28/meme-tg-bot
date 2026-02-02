package repo

import (
	"memetgbot/models"

	"gorm.io/gorm"
)

type ForwardMode struct {
	db *gorm.DB
}

func NewForwardModeRepo(db *gorm.DB) *ForwardMode {
	return &ForwardMode{db: db}
}

func (r *ForwardMode) Enable(userID int64, chatId int64) error {
	fm := models.ForwardMode{
		UserID:    userID,
		IsEnabled: true,
		ChatID:    chatId,
	}
	return r.db.Model(&models.ForwardMode{}).
		Where("user_id = ?", userID).
		Assign(fm).
		FirstOrCreate(&fm).
		Error
}

func (r *ForwardMode) Disable(userID int64) error {
	return r.db.
		Model(&models.ForwardMode{}).
		Where("user_id = ?", userID).
		Update("is_enabled", false).
		Error
}

func (r *ForwardMode) Get(userID int64) (models.ForwardMode, error) {
	var result models.ForwardMode

	resp := r.db.Find(&models.ForwardMode{}, &models.ForwardMode{UserID: userID}).Scan(&result)

	return result, resp.Error
}

func (r *ForwardMode) IsEnabled(userID int64) (bool, error) {
	fm, err := r.Get(userID)
	if err == nil {
		return fm.IsEnabled, nil
	}
	return false, err
}
