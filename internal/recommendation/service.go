package recommendation

import (
	"app/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RecommendationService struct {
	recommendationRepo *Repository
	redisClient        *redis.Client
	log                *logger.Logger
}

func NewRecommendationService(recommendationRepo *Repository, redisClient *redis.Client, log *logger.Logger) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		redisClient:        redisClient,
		log:                log,
	}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, userID uint) (recommendations []Recommendation, err error) {
	cacheKey := fmt.Sprintf("recommendations:user:%d", userID)

	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		if err = json.Unmarshal([]byte(val), &recommendations); err == nil {
			return recommendations, nil
		}
	}

	// Simulate a delay for the database call
	time.Sleep(3 * time.Second)
	recommendations, err = s.recommendationRepo.GetAllRecommendations()
	if err != nil {
		s.log.WithError(err).Error("Failed to get recommendations from database")
		return nil, err
	}

	cacheValue, err := json.Marshal(recommendations)
	if err != nil {
		s.log.WithError(err).Error("Failed to marshal recommendations for caching")
		return nil, err
	}

	err = s.redisClient.Set(ctx, cacheKey, cacheValue, 10*time.Minute).Err()
	if err != nil {
		s.log.WithError(err).Error("Failed to set recommendations in cache")
		return nil, err
	}

	return recommendations, nil
}
