package main

import (
	"Calsora/internal/config"
	"Calsora/internal/db"
	"Calsora/internal/handlers"
	"Calsora/internal/intelligence/nutrition"
	"Calsora/internal/repository"
	"Calsora/internal/server"
	"Calsora/internal/services"
	"Calsora/pkg/ratelimiter"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	conf := config.GetConfig()

	gin.SetMode(conf.GinMode)

	connDB, err := db.ConncetPostgres(conf.PostgresURL)
	if err != nil {
		log.Fatalf("Failed to cnnect to DB: %v\n", err)
	}
	if err := connDB.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	defer connDB.Close()
	log.Println("OK Connected to DB")

	redisClient, ok := ratelimiter.ConnectRedis()
	if ok {
		log.Println("OK Connected to Redis")
	}
	defer redisClient.Close()

	rateLimiterService := ratelimiter.NewRateLimiter(redisClient)

	userRepo := repository.NewUserRepository(connDB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	subRepo := repository.NewSubscriptionsRepository(connDB)
	subService := services.NewSubService(subRepo)
	subHandler := handlers.NewSubHandler(subService)

	authRepo := repository.NewAuthRepository(connDB)
	authService := services.NewAuthService(authRepo, userService, subService)
	authHandler := handlers.NewAuthHandler(authService)

	nutritionAI := nutrition.NewGPTClient()
	mealRepo := repository.NewMealRepository(connDB)
	mealService := services.NewMealService(mealRepo, nutritionAI)
	mealHandler := handlers.NewMealHandler(mealService)

	router := server.RunServer(userHandler, subHandler, authHandler, mealHandler, rateLimiterService)
	s := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: router,
	}

	go func() {
		log.Println("Starting server")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
	log.Println("Successful server stop")
}
