package main

import (
	"log"

	"fsa-boilerplate/backend/internal/api"
	"fsa-boilerplate/backend/internal/config"
	"fsa-boilerplate/backend/internal/dal"
)

func main() {
	cfg := config.Load()

	database, err := dal.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := dal.RunMigrations(database); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	r := api.New(database)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
