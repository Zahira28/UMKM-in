package main

import (
	"fmt"
	"log"

	_ "backend/docs"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title UMKM-in Backend API
// @version 1.0
// @description API backend service untuk aplikasi UMKM-in (Senior Project TI UGM)
// @contact.name Tim Senior Project UMKM-in
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT dengan format: Bearer <token>
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
		if err := database.SeedCategories(db); err != nil {
			log.Printf("Warning: Database category seeding failed: %v\n", err)
		}
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Swagger UI route
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html", fiber.StatusMovedPermanently)
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "UMKM-in backend service is running!",
		})
	})
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "UMKM-in backend service is running!",
		})
	})

	// Dependency injection for modules
	if db != nil {
		userRepo := repository.NewUserRepository(db)
		authService := service.NewAuthService(userRepo, cfg)
		authHandler := handler.NewAuthHandler(authService)

		categoryRepo := repository.NewCategoryRepository(db)
		categoryService := service.NewCategoryService(categoryRepo)
		categoryHandler := handler.NewCategoryHandler(categoryService)

		productRepo := repository.NewProductRepository(db)
		productService := service.NewProductService(productRepo, categoryRepo)
		productHandler := handler.NewProductHandler(productService)

		handler.SetupRoutes(app, authHandler, categoryHandler, productHandler, cfg.JWTSecret)
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server is running on http://localhost%s\n", addr)
	log.Printf("Swagger documentation is accessible at http://localhost%s/swagger\n", addr)
	log.Fatal(app.Listen(addr))
}