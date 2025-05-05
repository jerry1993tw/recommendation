package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(host, port, password string, db int) *redis.Client {
	addr := fmt.Sprintf("%s:%s", host, port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
