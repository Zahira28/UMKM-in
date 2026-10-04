package handler

import (
	"strconv"

	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/apperror"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	rawUserID := c.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.HandleError(c, apperror.Unauthorized("Autentikasi diperlukan"))
	}

	var req service.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format data request tidak valid"))
	}

	res, err := h.productService.CreateProduct(userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "Produk berhasil ditambahkan ke etalase", res)
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	categorySlug := c.Query("category", "")
	search := c.Query("search", "")
	sort := c.Query("sort", "latest")
	rawSellerID := c.Query("user_id", "")

	var sellerUUID *uuid.UUID
	if rawSellerID != "" {
		if parsed, err := uuid.Parse(rawSellerID); err == nil {
			sellerUUID = &parsed
		}
	}

	filter := repository.ProductFilter{
		CategorySlug: categorySlug,
		Search:       search,
		Sort:         sort,
		UserID:       sellerUUID,
		Page:         page,
		Limit:        limit,
	}

	var currentUserID *uuid.UUID
	if rawUserID := c.Locals("userID"); rawUserID != nil {
		if uid, ok := rawUserID.(uuid.UUID); ok && uid != uuid.Nil {
			currentUserID = &uid
		}
	}

	products, meta, err := h.productService.GetProducts(filter, currentUserID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.SuccessWithMeta(c, fiber.StatusOK, "Feed produk berhasil diambil", products, meta)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	rawID := c.Params("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID produk tidak valid"))
	}

	var currentUserID *uuid.UUID
	if rawUserID := c.Locals("userID"); rawUserID != nil {
		if uid, ok := rawUserID.(uuid.UUID); ok && uid != uuid.Nil {
			currentUserID = &uid
		}
	}

	product, err := h.productService.GetProductByID(id, currentUserID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Detail produk berhasil diambil", product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	rawID := c.Params("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID produk tidak valid"))
	}

	rawUserID := c.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.HandleError(c, apperror.Unauthorized("Autentikasi diperlukan"))
	}

	var req service.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format data request tidak valid"))
	}

	updatedProduct, err := h.productService.UpdateProduct(id, userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Produk berhasil diperbarui", updatedProduct)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	rawID := c.Params("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID produk tidak valid"))
	}

	rawUserID := c.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.HandleError(c, apperror.Unauthorized("Autentikasi diperlukan"))
	}

	if err := h.productService.DeleteProduct(id, userID); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Produk berhasil dihapus dari etalase", nil)
}
