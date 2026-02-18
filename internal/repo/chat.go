package repo

import (
	"errors"
	"memetgbot/model"

	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (chatRepo *ChatRepo) TryGet(chatId int64) (*model.Chat, bool, error) {
	var result model.Chat
	err := chatRepo.db.Where("telegram_id = ?", chatId).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return &result, true, nil
}

func (chatRepo *ChatRepo) Upsert(chat *model.Chat) error {
	return chatRepo.db.Where("telegram_id = ?", chat.TelegramID).
		Assign(chat).
		FirstOrCreate(chat).
		Error
}
