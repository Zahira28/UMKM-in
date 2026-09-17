package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/database"
)

func main() {
	log.Println("Starting database migration process...")

	cfg := config.LoadConfig()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Database auto-migration error: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
