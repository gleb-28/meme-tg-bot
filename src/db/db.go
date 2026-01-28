package db

import (
	"log"
	"memetgbot/models"
	"memetgbot/src/core/config"
	l "memetgbot/src/core/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(config *config.DatabaseConfig, logger l.Logger) *gorm.DB {
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(config.ToDSN()), &gorm.Config{})
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
