package main

import (
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize database connection and auto-migration
	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v. Ensure PostgreSQL is running.\n", err)
	} else {
		if err := database.AutoMigrate(db); err != nil {
			log.Printf("Warning: Database auto-migration failed: %v\n", err)
		}
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "UMKM-in backend service is running!",
		})
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(app.Listen(addr))
}