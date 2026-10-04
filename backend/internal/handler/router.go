package handler

import (
	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	authHandler *AuthHandler,
	categoryHandler *CategoryHandler,
	productHandler *ProductHandler,
	commentHandler *CommentHandler,
	likeHandler *LikeHandler,
	jwtSecret string,
) {
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/resend-otp", authHandler.ResendOTP)
	auth.Post("/login", authHandler.Login)
	auth.Post("/google", authHandler.GoogleAuth)

	// Protected profile routes
	authProtected := auth.Group("", middleware.Protected(jwtSecret))
	authProtected.Get("/profile", authHandler.GetProfile)
	authProtected.Put("/profile", authHandler.UpdateProfile)
	authProtected.Get("/me", authHandler.GetProfile)
	authProtected.Put("/me", authHandler.UpdateProfile)

	// Category routes
	categories := api.Group("/categories")
	categories.Get("", categoryHandler.GetCategories)

	// Product routes (Public / Optional Auth for feed and detail)
	products := api.Group("/products")
	products.Get("", middleware.OptionalAuth(jwtSecret), productHandler.GetProducts)
	products.Get("/:id", middleware.OptionalAuth(jwtSecret), productHandler.GetProductByID)

	// Product comments (Public)
	products.Get("/:id/comments", commentHandler.GetComments)

	// Protected product routes (Create, Update, Delete)
	productsProtected := products.Group("", middleware.Protected(jwtSecret))
	productsProtected.Post("", productHandler.CreateProduct)
	productsProtected.Put("/:id", productHandler.UpdateProduct)
	productsProtected.Delete("/:id", productHandler.DeleteProduct)

	// Comments & Likes on products (Protected)
	productsProtected.Post("/:id/comments", commentHandler.CreateComment)
	productsProtected.Post("/:id/like", likeHandler.ToggleLike)
	productsProtected.Post("/:id/react", likeHandler.ToggleLike)

	// Comment management (Protected - Edit & Delete)
	commentsProtected := api.Group("/comments", middleware.Protected(jwtSecret))
	commentsProtected.Put("/:id", commentHandler.UpdateComment)
	commentsProtected.Delete("/:id", commentHandler.DeleteComment)
}
