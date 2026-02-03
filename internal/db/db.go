package db

import (
	"fmt"
	"log"
	"memetgbot/internal/core/config"
	"memetgbot/internal/core/logger"
	"memetgbot/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func MustDB(config *config.AppConfig, logger logger.AppLogger) *gorm.DB {
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(sqlite.Open(config.Database.FileName), &gorm.Config{})
		if err == nil {
			logger.Info("Successfully connected to the database!")
			break
		}

		logger.Error(fmt.Sprintf("Failed to connect to database (attempt %d/%d): %v. Retrying in %v...\n", i+1, maxRetries, err, retryInterval))
		time.Sleep(retryInterval)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	err = db.AutoMigrate(&models.Chat{})
	if err != nil {
		log.Fatal("Failed to AutoMigrate Chat: " + err.Error())
	}

	err = db.AutoMigrate(&models.ForwardMode{})
	if err != nil {
		log.Fatal("Failed to AutoMigrate ForwardMode: " + err.Error())
	}

	return db
}
