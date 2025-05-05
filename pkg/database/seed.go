package database

import (
	"gorm.io/gorm"

	"app/internal/recommendation"
)

func Seed(db *gorm.DB) error {
	db.AutoMigrate(&recommendation.Recommendation{})

	var count int64
	db.Model(&recommendation.Recommendation{}).Count(&count)
	if count > 0 {
		return nil
	}

	recommendations := []recommendation.Recommendation{
		{Title: "Recommendation 1"},
		{Title: "Recommendation 2"},
		{Title: "Recommendation 3"},
		{Title: "Recommendation 4"},
		{Title: "Recommendation 5"},
	}

	return db.Create(&recommendations).Error
}
