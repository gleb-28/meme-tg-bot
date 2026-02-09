package model

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	TelegramID int64  `gorm:"uniqueIndex"`
	Name       string `gorm:"not null"`
	IsGroup    bool   `gorm:"not null;default:false"`
}

type ForwardMode struct {
	gorm.Model
	UserID    int64 `gorm:"uniqueIndex"`
	ChatID    int64
	IsEnabled bool `gorm:"not null;default:false"`

	Chat Chat `gorm:"foreignKey:UserID;references:ID"`
}
