package database

import (
	"app/internal/recommendation"
	"app/internal/user"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&recommendation.Recommendation{},
	)
}
