package handler

import (
	"backend/internal/service"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// GetCategories godoc
// @Summary Ambil seluruh daftar kategori produk
// @Description Mengambil semua kategori produk yang tersedia untuk filter dropdown di frontend.
// @Tags Categories
// @Produce json
// @Success 200 {object} response.APIResponse
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Daftar kategori berhasil diambil", categories)
}
