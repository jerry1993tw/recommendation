package recommendation

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RecommendationService struct {
	recommendationRepo *Repository
	redisClient        *redis.Client
}

func NewRecommendationService(recommendationRepo *Repository, redisClient *redis.Client) *RecommendationService {
	return &RecommendationService{
		recommendationRepo: recommendationRepo,
		redisClient:        redisClient,
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
		return nil, err
	}

	cacheValue, err := json.Marshal(recommendations)
	if err != nil {
		return nil, err
	}

	err = s.redisClient.Set(ctx, cacheKey, cacheValue, 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return recommendations, nil
}
