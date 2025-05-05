package recommendation

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllRecommendations() ([]Recommendation, error) {
	var recommendations []Recommendation
	err := r.db.Find(&recommendations).Error
	return recommendations, err
}
