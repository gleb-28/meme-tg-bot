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

func (chatRepo *ChatRepo) Upsert(chat *models.Chat) error {
	return chatRepo.db.Where("telegram_id = ?", chat.TelegramID).
		Assign(chat).
		FirstOrCreate(chat).
		Error
}

func (chatRepo *ChatRepo) Get(telegramID int64) (models.Chat, error) {
	var result models.Chat
	resp := chatRepo.db.Find(&models.Chat{}, &models.Chat{TelegramID: telegramID}).Scan(&result)

	return result, resp.Error
}
