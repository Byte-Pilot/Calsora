package main

import (
	"Calsora/internal/config"
	"Calsora/internal/db"
	"context"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	conf := config.GetConfig()
	conn, err := db.ConncetPostgres(conf.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("OK")
}
