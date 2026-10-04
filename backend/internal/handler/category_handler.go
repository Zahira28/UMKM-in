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

func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetCategories()
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Daftar kategori berhasil diambil", categories)
}
