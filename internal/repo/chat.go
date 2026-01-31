package repo

import (
	"memetgbot/internal/db"
	"memetgbot/models"
)

type ChatRepo struct{}

func (chatRepo *ChatRepo) Add(telegramID int64) error {
	resp := db.DB.Model(&models.Chat{}).Create(&models.Chat{
		TelegramID: telegramID,
	})

	return resp.Error
}

func (chatRepo *ChatRepo) Get(telegramID int64) (models.Chat, error) {
	var result models.Chat
	resp := db.DB.Find(&models.Chat{}, &models.Chat{TelegramID: telegramID}).Scan(&result)

	return result, resp.Error
}

var Chat = ChatRepo{}
