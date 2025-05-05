package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/auth"
	"app/internal/config"
	"app/internal/middleware"
	"app/internal/recommendation"
	"app/internal/user"
	"app/pkg/email"
	"app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	dbConn      *gorm.DB
	redisClient *redis.Client
)

func Start(db *gorm.DB, redis *redis.Client) {
	dbConn = db
	redisClient = redis

	router := InitRouter()

	s := &http.Server{
		Addr:         ":" + config.New().Server.Port,
		ReadTimeout:  time.Duration(config.New().Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.New().Server.WriteTimeout) * time.Second,
		Handler:      router,
	}

	// Start server
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Println("Database connection close:", err)
		}
	}

	if err := redisClient.Close(); err != nil {
		log.Println("Redis connection close:", err)
	}
}

// InitRouter initializes routes
func InitRouter() *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	log := logger.New()

	userRepo := user.NewRepository(dbConn)
	emailService := &email.MockEmailService{}
	authService := auth.NewService(userRepo, emailService, log)
	authHandler := auth.NewHandler(authService, log)

	recommendationRepo := recommendation.NewRepository(dbConn)
	recommendationService := recommendation.NewRecommendationService(recommendationRepo, redisClient, log)
	recommendationHandler := recommendation.NewHandler(recommendationService, log)

	e.POST("/register", authHandler.Register)
	e.POST("/verify", authHandler.VerifyEmail)
	e.POST("/login", authHandler.Login)
	e.GET("/recommendations", middleware.AuthMiddleware(), recommendationHandler.GetRecommendations)

	return e
}
