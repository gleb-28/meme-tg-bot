package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	TelegramID uint32 `gorm:"unique"`
}
