package db

import (
	"log"
	c "memetgbot/internal/core/config"
	l "memetgbot/internal/core/logger"
	"memetgbot/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var config = c.Config
	var logger = l.Logger
	var err error
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(sqlite.Open(config.Database.FileName), &gorm.Config{})
		if err == nil {
			logger.Info("Successfully connected to the database!")
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...\n", i+1, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	err = DB.AutoMigrate(&models.Chat{})
	if err != nil {
		panic("Failed to AutoMigrate: " + err.Error())
	}
}
