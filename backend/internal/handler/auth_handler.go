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

// Register godoc
// @Summary Register pengguna baru
// @Description Mendaftar akun UMKM-in menggunakan nama lengkap, email, dan kata sandi. Kode OTP akan dikirimkan ke email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Data Registrasi Akun"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Router /auth/register [post]
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

// VerifyOTP godoc
// @Summary Verifikasi OTP Email
// @Description Memverifikasi kode OTP 6 digit yang dikirimkan ke email untuk mengaktifkan akun.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.VerifyOTPRequest true "Data Verifikasi OTP"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/verify-otp [post]
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

// ResendOTP godoc
// @Summary Kirim ulang kode OTP
// @Description Mengirimkan kode OTP baru ke email yang terdaftar jika kode sebelumnya kedaluwarsa.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.ResendOTPRequest true "Email pengguna"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/resend-otp [post]
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

// Login godoc
// @Summary Login akun manual
// @Description Masuk ke akun menggunakan email atau username terdaftar beserta kata sandi.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Kredensial Login"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/login [post]
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

// GoogleAuth godoc
// @Summary Login / Daftar via Google OAuth 2.0
// @Description Melakukan autentikasi menggunakan Google ID token dari frontend Google Identity Services.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service.GoogleAuthRequest true "Google ID Token"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /auth/google [post]
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

// GetProfile godoc
// @Summary Ambil profil pengguna yang sedang login
// @Description Mendapatkan data profil akun pemilik token JWT aktif.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/profile [get]
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

// UpdateProfile godoc
// @Summary Lengkapi / Perbarui profil pengguna
// @Description Memperbarui informasi profil toko atau melengkapi data setelah onboarding.
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.UpdateProfileRequest true "Data Pembaruan Profil"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /auth/profile [put]
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

	updatedUser, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, "Profil berhasil diperbarui", updatedUser)
}
