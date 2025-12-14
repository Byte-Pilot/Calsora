package main

import (
	"Calsora/internal/config"
	"Calsora/internal/db"
	"Calsora/internal/handlers"
	"Calsora/internal/repository"
	"Calsora/internal/server"
	"Calsora/internal/services"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(".env file not found")
	}

	conf := config.GetConfig()

	gin.SetMode(conf.GinMode)

	conn, err := db.ConncetPostgres(conf.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("OK")

	userRepo := repository.NewUserRepository(conn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	mealRepo := repository.NewMealRepository(conn)
	mealService := services.NewMealService(mealRepo)
	mealHandler := handlers.NewMealHandler(mealService)

	s := server.RunServer(userHandler, mealHandler)
	if err := s.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
