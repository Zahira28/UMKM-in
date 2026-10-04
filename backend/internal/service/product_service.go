package service

import (
	"math"
	"strings"
	"time"

	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperror"
	"backend/pkg/response"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	CategoryID     uint    `json:"category_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Unit           string  `json:"unit"`
	RawImageURL    string  `json:"raw_image_url"`
	StudioImageURL string  `json:"studio_image_url"`
	IsAvailable    *bool   `json:"is_available"`
}

type UpdateProductRequest struct {
	CategoryID     *uint    `json:"category_id"`
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	Price          *float64 `json:"price"`
	Unit           *string  `json:"unit"`
	RawImageURL    *string  `json:"raw_image_url"`
	StudioImageURL *string  `json:"studio_image_url"`
	IsAvailable    *bool    `json:"is_available"`
}

type ProductCreatedResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductSellerResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	AvatarURL   string    `json:"avatar_url"`
	City        string    `json:"city"`
	PhoneNumber string    `json:"phone_number"`
}

type ProductResponse struct {
	ID             uuid.UUID              `json:"id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Price          float64                `json:"price"`
	Unit           string                 `json:"unit"`
	RawImageURL    string                 `json:"raw_image_url"`
	StudioImageURL string                 `json:"studio_image_url"`
	IsAvailable    bool                   `json:"is_available"`
	CreatedAt      time.Time              `json:"created_at"`
	LikesCount     int64                  `json:"likes_count"`
	CommentsCount  int64                  `json:"comments_count"`
	IsLiked        bool                   `json:"is_liked"`
	Category       *CategoryResponse      `json:"category"`
	Seller         *ProductSellerResponse `json:"seller"`
}

type ProductService interface {
	CreateProduct(userID uuid.UUID, req CreateProductRequest) (*ProductCreatedResponse, error)
	GetProducts(filter repository.ProductFilter, currentUserID *uuid.UUID) ([]ProductResponse, *response.PaginationMeta, error)
	GetProductByID(id uuid.UUID, currentUserID *uuid.UUID) (*ProductResponse, error)
	UpdateProduct(id uuid.UUID, currentUserID uuid.UUID, req UpdateProductRequest) (*ProductResponse, error)
	DeleteProduct(id uuid.UUID, currentUserID uuid.UUID) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) CreateProduct(userID uuid.UUID, req CreateProductRequest) (*ProductCreatedResponse, error) {
	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		return nil, apperror.BadRequest("Nama produk minimal 3 karakter")
	}

	if req.Price < 0 {
		return nil, apperror.BadRequest("Harga produk tidak boleh negatif")
	}

	unit := strings.TrimSpace(req.Unit)
	if unit == "" {
		return nil, apperror.BadRequest("Satuan produk wajib diisi (contoh: pcs, box, porsi)")
	}

	if req.CategoryID == 0 {
		return nil, apperror.BadRequest("Kategori produk wajib dipilih")
	}

	category, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, apperror.Internal("Gagal memvalidasi kategori: " + err.Error())
	}
	if category == nil {
		return nil, apperror.BadRequest("Kategori produk tidak ditemukan atau tidak valid")
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	product := &model.Product{
		UserID:         userID,
		CategoryID:     req.CategoryID,
		Name:           name,
		Description:    strings.TrimSpace(req.Description),
		Price:          req.Price,
		Unit:           unit,
		RawImageURL:    strings.TrimSpace(req.RawImageURL),
		StudioImageURL: strings.TrimSpace(req.StudioImageURL),
		IsAvailable:    isAvailable,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, apperror.Internal("Gagal menambahkan produk ke etalase: " + err.Error())
	}

	return &ProductCreatedResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		IsAvailable: product.IsAvailable,
		CreatedAt:   product.CreatedAt,
	}, nil
}

func (s *productService) GetProducts(filter repository.ProductFilter, currentUserID *uuid.UUID) ([]ProductResponse, *response.PaginationMeta, error) {
	products, totalItems, err := s.productRepo.FindAll(filter)
	if err != nil {
		return nil, nil, apperror.Internal("Gagal mengambil daftar produk: " + err.Error())
	}

	productIDs := make([]uuid.UUID, 0, len(products))
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}

	likesCountMap, commentsCountMap, isLikedMap, err := s.productRepo.GetProductStats(productIDs, currentUserID)
	if err != nil {
		return nil, nil, apperror.Internal("Gagal mengambil statistik produk: " + err.Error())
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(limit)))
	}

	meta := &response.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	result := make([]ProductResponse, 0, len(products))
	for _, p := range products {
		res := s.mapProductToResponse(&p, likesCountMap[p.ID], commentsCountMap[p.ID], isLikedMap[p.ID])
		result = append(result, res)
	}

	return result, meta, nil
}

func (s *productService) GetProductByID(id uuid.UUID, currentUserID *uuid.UUID) (*ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("Gagal mengambil detail produk: " + err.Error())
	}
	if product == nil {
		return nil, apperror.NotFound("Produk tidak ditemukan")
	}

	likesCountMap, commentsCountMap, isLikedMap, err := s.productRepo.GetProductStats([]uuid.UUID{product.ID}, currentUserID)
	if err != nil {
		return nil, apperror.Internal("Gagal mengambil statistik produk: " + err.Error())
	}

	res := s.mapProductToResponse(product, likesCountMap[product.ID], commentsCountMap[product.ID], isLikedMap[product.ID])
	return &res, nil
}

func (s *productService) UpdateProduct(id uuid.UUID, currentUserID uuid.UUID, req UpdateProductRequest) (*ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("Gagal memuat produk: " + err.Error())
	}
	if product == nil {
		return nil, apperror.NotFound("Produk tidak ditemukan")
	}

	// Ownership check
	if product.UserID != currentUserID {
		return nil, apperror.Forbidden("Anda tidak memiliki akses untuk mengubah produk ini")
	}

	if req.CategoryID != nil {
		cat, err := s.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, apperror.Internal("Gagal memvalidasi kategori: " + err.Error())
		}
		if cat == nil {
			return nil, apperror.BadRequest("Kategori tidak valid atau tidak ditemukan")
		}
		product.CategoryID = *req.CategoryID
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 {
			return nil, apperror.BadRequest("Nama produk minimal 3 karakter")
		}
		product.Name = name
	}

	if req.Description != nil {
		product.Description = strings.TrimSpace(*req.Description)
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return nil, apperror.BadRequest("Harga produk tidak boleh negatif")
		}
		product.Price = *req.Price
	}

	if req.Unit != nil {
		unit := strings.TrimSpace(*req.Unit)
		if unit == "" {
			return nil, apperror.BadRequest("Satuan produk tidak boleh kosong")
		}
		product.Unit = unit
	}

	if req.RawImageURL != nil {
		product.RawImageURL = strings.TrimSpace(*req.RawImageURL)
	}

	if req.StudioImageURL != nil {
		product.StudioImageURL = strings.TrimSpace(*req.StudioImageURL)
	}

	if req.IsAvailable != nil {
		product.IsAvailable = *req.IsAvailable
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, apperror.Internal("Gagal memperbarui data produk: " + err.Error())
	}

	// Reload updated product with relations
	updatedProduct, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, apperror.Internal("Gagal memuat data produk terbaru: " + err.Error())
	}

	likesCountMap, commentsCountMap, isLikedMap, _ := s.productRepo.GetProductStats([]uuid.UUID{id}, &currentUserID)
	res := s.mapProductToResponse(updatedProduct, likesCountMap[id], commentsCountMap[id], isLikedMap[id])
	return &res, nil
}

func (s *productService) DeleteProduct(id uuid.UUID, currentUserID uuid.UUID) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return apperror.Internal("Gagal memuat data produk: " + err.Error())
	}
	if product == nil {
		return apperror.NotFound("Produk tidak ditemukan")
	}

	// Ownership check
	if product.UserID != currentUserID {
		return apperror.Forbidden("Anda tidak memiliki akses untuk menghapus produk ini")
	}

	if err := s.productRepo.Delete(id); err != nil {
		return apperror.Internal("Gagal menghapus produk: " + err.Error())
	}

	return nil
}

func (s *productService) mapProductToResponse(p *model.Product, likesCount int64, commentsCount int64, isLiked bool) ProductResponse {
	var categoryResp *CategoryResponse
	if p.Category != nil {
		categoryResp = &CategoryResponse{
			ID:   p.Category.ID,
			Name: p.Category.Name,
			Slug: p.Category.Slug,
		}
	}

	var sellerResp *ProductSellerResponse
	if p.User != nil {
		sellerResp = &ProductSellerResponse{
			ID:          p.User.ID,
			Username:    p.User.Username,
			FullName:    p.User.FullName,
			AvatarURL:   p.User.AvatarURL,
			City:        p.User.City,
			PhoneNumber: p.User.PhoneNumber,
		}
	}

	return ProductResponse{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Price:          p.Price,
		Unit:           p.Unit,
		RawImageURL:    p.RawImageURL,
		StudioImageURL: p.StudioImageURL,
		IsAvailable:    p.IsAvailable,
		CreatedAt:      p.CreatedAt,
		LikesCount:     likesCount,
		CommentsCount:  commentsCount,
		IsLiked:        isLiked,
		Category:       categoryResp,
		Seller:         sellerResp,
	}
}
