package db

import (
	"log"
	c "memetgbot/internal/core/config"
	l "memetgbot/internal/core/logger"
	"memetgbot/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var config = c.Config
	var logger = l.Logger
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(config.Database.ToDSN()), &gorm.Config{})
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

	err = db.AutoMigrate(&models.Chat{})
	if err != nil {
		panic("Failed to AutoMigrate: " + err.Error())
	}

	return db
}
