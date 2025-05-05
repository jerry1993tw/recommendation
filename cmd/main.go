package main

import (
	"context"
	"fmt"

	"app/internal/config"
	"app/pkg/cache"
	"app/pkg/database"
	"app/pkg/logger"
	"app/pkg/server"

	"github.com/google/uuid"
)

func main() {
	ctx := context.WithValue(context.TODO(), "requestId", uuid.New().String())
	log := logger.New().WithContext(ctx)

	log.Info("Connecting to database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.New().Database.User,
		config.New().Database.Password,
		config.New().Database.Host,
		config.New().Database.Port,
		config.New().Database.Name,
	)
	database.InitDB(dsn)
	dbConn := database.GetDB()

	log.Info("Running database migrations")
	if err := database.AutoMigrate(dbConn); err != nil {
		log.WithError(err).Error("Failed to run database migrations")
		return
	}

	log.Info("Initializing Redis client")
	redisClient := cache.NewRedisClient(
		config.New().Cache.Host,
		config.New().Cache.Port,
		config.New().Cache.Password,
		config.New().Cache.DB,
	)

	log.Info("Seeding database")
	if err := database.Seed(dbConn); err != nil {
		log.WithError(err).Error("Failed to seed database")
		return
	}

	log.Info("Starting server")
	server.Start(dbConn, redisClient)
}
