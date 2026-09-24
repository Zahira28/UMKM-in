package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"backend/pkg/apperror"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyOTPRequest struct {
	Email   string `json:"email"`
	OTPCode string `json:"otp_code"`
}

type ResendOTPRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"id_token"`
}

type UpdateProfileRequest struct {
	Username      string   `json:"username"`
	FullName      string   `json:"full_name"`
	PhoneNumber   string   `json:"phone_number"`
	Bio           string   `json:"bio"`
	AvatarURL     string   `json:"avatar_url"`
	City          string   `json:"city"`
	AddressDetail string   `json:"address_detail"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type AuthService interface {
	Register(req RegisterRequest) (*model.User, error)
	VerifyOTP(req VerifyOTPRequest) (*AuthResponse, error)
	ResendOTP(req ResendOTPRequest) error
	Login(req LoginRequest) (*AuthResponse, error)
	GoogleAuth(req GoogleAuthRequest) (*AuthResponse, error)
	GetProfile(userID uuid.UUID) (*model.User, error)
	UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*model.User, error)
}

type authService struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	googleClientID string
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo:       userRepo,
		jwtSecret:      cfg.JWTSecret,
		googleClientID: cfg.GoogleClientID,
	}
}

func (s *authService) Register(req RegisterRequest) (*model.User, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.FullName = strings.TrimSpace(req.FullName)

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return nil, apperror.BadRequest("Nama lengkap, email, dan kata sandi wajib diisi")
	}

	if len(req.Password) < 6 {
		return nil, apperror.BadRequest("Kata sandi minimal 6 karakter")
	}

	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperror.Internal("Gagal memeriksa ketersediaan email di database", err)
	}
	if existingUser != nil {
		return nil, apperror.Conflict("Email sudah terdaftar dalam sistem")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, apperror.Internal("Gagal mengenkripsi kata sandi", err)
	}

	otpCode, err := utils.GenerateOTP()
	if err != nil {
		return nil, apperror.Internal("Gagal menghasilkan kode verifikasi", err)
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	user := &model.User{
		FullName:         req.FullName,
		Email:            req.Email,
		PasswordHash:     hashedPassword,
		IsVerified:       false,
		VerificationCode: otpCode,
		CodeExpiresAt:    &expiresAt,
		AICredits:        10,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, apperror.Internal("Gagal mendaftarkan pengguna ke basis data", err)
	}

	log.Printf("[DEV OTP] Verification code for %s is: %s (expires in 10 minutes)\n", user.Email, otpCode)

	return user, nil
}

func (s *authService) VerifyOTP(req VerifyOTPRequest) (*AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.OTPCode = strings.TrimSpace(req.OTPCode)

	if req.Email == "" || req.OTPCode == "" {
		return nil, apperror.BadRequest("Email dan kode OTP wajib diisi")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, apperror.Internal("Gagal mengambil data akun", err)
	}
	if user == nil {
		return nil, apperror.NotFound("Akun dengan email tersebut tidak ditemukan")
	}

	if user.IsVerified {
		token, err := utils.GenerateToken(user.ID, user.Email, user.Username, s.jwtSecret, 7*24*time.Hour)
		if err != nil {
			return nil, apperror.Internal("Gagal menerbitkan token sesi", err)
		}
		return &AuthResponse{Token: token, User: user}, nil
	}

	if user.VerificationCode != req.OTPCode {
		return nil, apperror.BadRequest("Kode OTP yang Anda masukkan salah")
	}

	if user.CodeExpiresAt == nil || time.Now().After(*user.CodeExpiresAt) {
		return nil, apperror.BadRequest("Kode OTP sudah kedaluwarsa, silakan minta kode baru")
	}

	user.IsVerified = true
	user.VerificationCode = ""
	user.CodeExpiresAt = nil

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperror.Internal("Gagal memperbarui status verifikasi akun", err)
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, s.jwtSecret, 7*24*time.Hour)
	if err != nil {
		return nil, apperror.Internal("Gagal menerbitkan token sesi", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) ResendOTP(req ResendOTPRequest) error {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" {
		return apperror.BadRequest("Email wajib diisi")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return apperror.Internal("Gagal memeriksa status akun", err)
	}
	if user == nil {
		return apperror.NotFound("Akun tidak ditemukan")
	}

	if user.IsVerified {
		return apperror.BadRequest("Akun ini sudah terverifikasi sebelumnya")
	}

	newOTP, err := utils.GenerateOTP()
	if err != nil {
		return apperror.Internal("Gagal menghasilkan kode verifikasi baru", err)
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	user.VerificationCode = newOTP
	user.CodeExpiresAt = &expiresAt

	if err := s.userRepo.Update(user); err != nil {
		return apperror.Internal("Gagal memperbarui kode verifikasi", err)
	}

	log.Printf("[DEV OTP RESEND] New verification code for %s is: %s\n", user.Email, newOTP)
	return nil
}

func (s *authService) Login(req LoginRequest) (*AuthResponse, error) {
	req.EmailOrUsername = strings.TrimSpace(req.EmailOrUsername)

	if req.EmailOrUsername == "" || req.Password == "" {
		return nil, apperror.BadRequest("Email/username dan kata sandi wajib diisi")
	}

	user, err := s.userRepo.FindByEmailOrUsername(req.EmailOrUsername)
	if err != nil {
		return nil, apperror.Internal("Gagal membaca data akun", err)
	}
	if user == nil {
		return nil, apperror.Unauthorized("Email/username atau kata sandi tidak valid")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, apperror.Unauthorized("Email/username atau kata sandi tidak valid")
	}

	if !user.IsVerified {
		return nil, apperror.Unauthorized("Akun Anda belum diverifikasi, silakan lakukan verifikasi kode OTP terlebih dahulu")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, s.jwtSecret, 7*24*time.Hour)
	if err != nil {
		return nil, apperror.Internal("Gagal menerbitkan token autentikasi", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

type googleTokenInfo struct {
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Error         string `json:"error_description"`
}

func (s *authService) GoogleAuth(req GoogleAuthRequest) (*AuthResponse, error) {
	if req.IDToken == "" {
		return nil, apperror.BadRequest("ID Token Google wajib disertakan")
	}

	verifyURL := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", req.IDToken)
	resp, err := http.Get(verifyURL)
	if err != nil {
		return nil, apperror.Internal("Gagal memverifikasi token ke server Google", err)
	}
	defer resp.Body.Close()

	var tokenInfo googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, apperror.Internal("Gagal membaca respon verifikasi Google", err)
	}

	if tokenInfo.Error != "" || tokenInfo.Email == "" {
		return nil, apperror.Unauthorized("Token Google tidak valid atau sudah kedaluwarsa")
	}

	// Validasi audience jika GOOGLE_CLIENT_ID diset di environment
	if s.googleClientID != "" && tokenInfo.Aud != s.googleClientID {
		return nil, apperror.Unauthorized("Token Google tidak ditujukan untuk aplikasi ini")
	}

	tokenInfo.Email = strings.ToLower(tokenInfo.Email)

	user, err := s.userRepo.FindByEmail(tokenInfo.Email)
	if err != nil {
		return nil, apperror.Internal("Gagal memeriksa akun Google di database", err)
	}

	if user == nil {
		baseUsername := strings.Split(tokenInfo.Email, "@")[0]
		user = &model.User{
			Email:      tokenInfo.Email,
			FullName:   tokenInfo.Name,
			Username:   baseUsername,
			AvatarURL:  tokenInfo.Picture,
			GoogleID:   tokenInfo.Sub,
			IsVerified: true,
			AICredits:  10,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, apperror.Internal("Gagal mendaftarkan akun Google baru", err)
		}
	} else {
		user.IsVerified = true
		if user.GoogleID == "" {
			user.GoogleID = tokenInfo.Sub
		}
		if user.AvatarURL == "" && tokenInfo.Picture != "" {
			user.AvatarURL = tokenInfo.Picture
		}
		_ = s.userRepo.Update(user)
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, s.jwtSecret, 7*24*time.Hour)
	if err != nil {
		return nil, apperror.Internal("Gagal menerbitkan token sesi", err)
	}

	return &AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) GetProfile(userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.Internal("Gagal mengambil data profil", err)
	}
	if user == nil {
		return nil, apperror.NotFound("Pengguna tidak ditemukan")
	}
	return user, nil
}

func (s *authService) UpdateProfile(userID uuid.UUID, req UpdateProfileRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, apperror.Internal("Gagal mencari data profil pengguna", err)
	}
	if user == nil {
		return nil, apperror.NotFound("Pengguna tidak ditemukan")
	}

	if req.Username != "" && strings.ToLower(req.Username) != strings.ToLower(user.Username) {
		existing, err := s.userRepo.FindByUsername(req.Username)
		if err != nil {
			return nil, apperror.Internal("Gagal memvalidasi ketersediaan username", err)
		}
		if existing != nil && existing.ID != user.ID {
			return nil, apperror.Conflict("Username ini sudah digunakan oleh akun lain")
		}
		user.Username = strings.ToLower(strings.TrimSpace(req.Username))
	}

	if req.FullName != "" {
		user.FullName = strings.TrimSpace(req.FullName)
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = strings.TrimSpace(req.PhoneNumber)
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.AddressDetail != "" {
		user.AddressDetail = req.AddressDetail
	}
	if req.Latitude != nil {
		user.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		user.Longitude = req.Longitude
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, apperror.Internal("Gagal memperbarui data profil ke database", err)
	}

	return user, nil
}
