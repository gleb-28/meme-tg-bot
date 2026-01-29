package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	TelegramID int64 `gorm:"unique"`
}
