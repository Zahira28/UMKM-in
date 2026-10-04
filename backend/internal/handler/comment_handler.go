package handler

import (
	"strconv"

	"backend/internal/service"
	"backend/pkg/apperror"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// GetComments godoc
// @Summary Ambil daftar komentar produk
// @Description Mengambil semua komentar pada produk beserta balasan bertingkat (nested replies).
// @Tags Comments
// @Produce json
// @Param id path string true "UUID Produk"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /products/{id}/comments [get]
func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	rawProductID := c.Params("id")
	productID, err := uuid.Parse(rawProductID)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID produk tidak valid"))
	}

	comments, err := h.commentService.GetComments(productID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Daftar komentar berhasil diambil", comments)
}

// CreateComment godoc
// @Summary Kirim komentar baru atau balasan
// @Description Menambahkan komentar baru ke produk atau membalas komentar yang sudah ada dengan menyertakan parent_id.
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Produk"
// @Param request body service.CreateCommentRequest true "Isi Komentar"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /products/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	rawProductID := c.Params("id")
	productID, err := uuid.Parse(rawProductID)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID produk tidak valid"))
	}

	rawUserID := c.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.HandleError(c, apperror.Unauthorized("Autentikasi diperlukan"))
	}

	var req service.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format data request tidak valid"))
	}

	res, err := h.commentService.CreateComment(productID, userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "Komentar berhasil dikirim", res)
}

// UpdateComment godoc
// @Summary Edit teks komentar
// @Description Memperbarui isi teks komentar yang pernah dibuat. Hanya dapat dilakukan oleh pembuat komentar itu sendiri.
// @Tags Comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Komentar"
// @Param request body service.UpdateCommentRequest true "Data Pembaruan Komentar"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	commentID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID komentar tidak valid"))
	}

	rawUserID := c.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.HandleError(c, apperror.Unauthorized("Autentikasi diperlukan"))
	}

	var req service.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format data request tidak valid"))
	}

	res, err := h.commentService.UpdateComment(commentID, userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Komentar berhasil diperbarui", res)
}

// DeleteComment godoc
// @Summary Hapus komentar
// @Description Menghapus komentar dari produk. Dapat dilakukan oleh pembuat komentar atau pemilik produk/toko.
// @Tags Comments
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Komentar"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	commentID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.HandleError(c, apperror.BadRequest("Format ID komentar tidak valid"))
	}

	rawUserID := c.Locals("userID")
	userID, ok := rawUserID.(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return response.HandleError(c, apperror.Unauthorized("Autentikasi diperlukan"))
	}

	if err := h.commentService.DeleteComment(commentID, userID); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Komentar berhasil dihapus", nil)
}
