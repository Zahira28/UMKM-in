package handler

import (
	"backend/internal/service"
	"backend/pkg/apperror"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LikeHandler struct {
	likeService service.LikeService
}

func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// ToggleLike godoc
// @Summary Toggle suka / reaction pada produk
// @Description Menyukai produk atau membatalkan suka (unlike) pada produk. Memerlukan autentikasi.
// @Tags Likes
// @Produce json
// @Security BearerAuth
// @Param id path string true "UUID Produk"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /products/{id}/like [post]
func (h *LikeHandler) ToggleLike(c *fiber.Ctx) error {
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

	res, err := h.likeService.ToggleLike(productID, userID)
	if err != nil {
		return response.HandleError(c, err)
	}

	message := "Berhasil menyukai produk"
	if !res.IsLiked {
		message = "Berhasil membatalkan suka pada produk"
	}

	return response.Success(c, fiber.StatusOK, message, res)
}
