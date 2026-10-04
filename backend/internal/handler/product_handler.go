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

// CreateProduct godoc
// @Summary Tambah produk baru ke etalase
// @Description Menambahkan data produk baru ke etalase toko. Memerlukan autentikasi seller.
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateProductRequest true "Data Produk Baru"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /products [post]
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

// GetProducts godoc
// @Summary Ambil linimasa / feed produk publik
// @Description Mengambil katalog produk dengan opsi filter kategori, pencarian teks, pengurutan harga/waktu, dan pagination.
// @Tags Products
// @Produce json
// @Param category query string false "Filter slug kategori (makanan, minuman, pakaian)"
// @Param search query string false "Kata kunci pencarian judul atau deskripsi"
// @Param sort query string false "Urutan sorting: latest (default), cheapest, priciest"
// @Param user_id query string false "UUID Penjual (filter produk milik seller tertentu)"
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah item per halaman (default: 10, max: 100)"
// @Success 200 {object} response.APIResponseWithMeta
// @Router /products [get]
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

// GetProductByID godoc
// @Summary Ambil detail produk berdasarkan ID
// @Description Mengambil rincian informasi lengkap sebuah produk beserta informasi toko penjual.
// @Tags Products
// @Produce json
// @Param id path string true "UUID Produk"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /products/{id} [get]
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

// UpdateProduct godoc
// @Summary Perbarui data produk
// @Description Mengubah rincian informasi produk di etalase. Hanya dapat dilakukan oleh pemilik produk.
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Produk"
// @Param request body service.UpdateProductRequest true "Data Produk Baru"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /products/{id} [put]
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

// DeleteProduct godoc
// @Summary Hapus produk dari etalase
// @Description Menghapus produk dari etalase toko secara permanen. Hanya dapat dilakukan oleh pemilik produk.
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Produk"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /products/{id} [delete]
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
