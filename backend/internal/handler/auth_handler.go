package handler

import (
	"backend/internal/service"
	"backend/pkg/apperror"
	"backend/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format JSON tidak valid", err.Error()))
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, "Pendaftaran berhasil. Silakan cek email Anda untuk kode verifikasi OTP.", fiber.Map{
		"email":              user.Email,
		"expires_in_seconds": 600,
	})
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req service.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format JSON tidak valid", err.Error()))
	}

	authResp, err := h.authService.VerifyOTP(req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Verifikasi email berhasil, akun Anda telah aktif", authResp)
}

func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
	var req service.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format JSON tidak valid", err.Error()))
	}

	if err := h.authService.ResendOTP(req); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Kode OTP baru telah berhasil dikirim ke email Anda", fiber.Map{
		"cooldown_seconds": 60,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format JSON tidak valid", err.Error()))
	}

	authResp, err := h.authService.Login(req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Login berhasil", authResp)
}

func (h *AuthHandler) GoogleAuth(c *fiber.Ctx) error {
	var req service.GoogleAuthRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format JSON tidak valid", err.Error()))
	}

	authResp, err := h.authService.GoogleAuth(req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Autentikasi Google berhasil", authResp)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return response.HandleError(c, apperror.Unauthorized("Sesi login tidak sah"))
	}

	user, err := h.authService.GetProfile(userID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Data profil berhasil diambil", user)
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return response.HandleError(c, apperror.Unauthorized("Sesi login tidak sah"))
	}

	var req service.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.HandleError(c, apperror.BadRequest("Format JSON tidak valid", err.Error()))
	}

	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Profil berhasil diperbarui", user)
}
