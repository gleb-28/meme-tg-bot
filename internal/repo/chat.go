package repo

import (
	"memetgbot/model"

	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(db *gorm.DB) *ChatRepo {
	return &ChatRepo{db: db}
}

func (chatRepo *ChatRepo) Upsert(chat *model.Chat) error {
	return chatRepo.db.Where("telegram_id = ?", chat.TelegramID).
		Assign(chat).
		FirstOrCreate(chat).
		Error
}

func (chatRepo *ChatRepo) Get(chatId int64) (model.Chat, error) {
	var result model.Chat
	resp := chatRepo.db.Find(&model.Chat{}, &model.Chat{TelegramID: chatId}).Scan(&result)

	return result, resp.Error
}
