package repo

import (
	"memetgbot/models"

	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (chatRepo *ChatRepo) Add(telegramID int64) error {
	resp := chatRepo.db.Model(&models.Chat{}).Create(&models.Chat{
		TelegramID: telegramID,
	})

	return resp.Error
}

func (chatRepo *ChatRepo) Get(telegramID int64) (models.Chat, error) {
	var result models.Chat
	resp := chatRepo.db.Find(&models.Chat{}, &models.Chat{TelegramID: telegramID}).Scan(&result)

	return result, resp.Error
}
