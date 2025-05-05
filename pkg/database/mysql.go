package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(dsn string) (*gorm.DB, error) {
	var err error
	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Printf("[database] Connected to MySQL successfully (attempt %d)\n", i)
			return db, nil
		}

		log.Printf("[database] Failed to connect to MySQL (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to MySQL after %d attempts: %w", maxRetries, err)
}

func GetDB() *gorm.DB {
	return db
}
