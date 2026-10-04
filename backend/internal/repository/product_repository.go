package repository

import (
	"errors"
	"strings"

	"backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductFilter struct {
	CategorySlug string
	CategoryID   uint
	Search       string
	Sort         string
	UserID       *uuid.UUID
	Page         int
	Limit        int
}

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id uuid.UUID) (*model.Product, error)
	FindAll(filter ProductFilter) ([]model.Product, int64, error)
	Update(product *model.Product) error
	Delete(id uuid.UUID) error
	GetProductStats(productIDs []uuid.UUID, currentUserID *uuid.UUID) (likesCount map[uuid.UUID]int64, commentsCount map[uuid.UUID]int64, isLiked map[uuid.UUID]bool, err error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Preload("User").First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(filter ProductFilter) ([]model.Product, int64, error) {
	query := r.db.Model(&model.Product{}).Preload("Category").Preload("User")

	// Filter by Category Slug
	if strings.TrimSpace(filter.CategorySlug) != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("LOWER(categories.slug) = LOWER(?)", strings.TrimSpace(filter.CategorySlug))
	} else if filter.CategoryID > 0 {
		query = query.Where("products.category_id = ?", filter.CategoryID)
	}

	// Filter by UserID (Seller)
	if filter.UserID != nil && *filter.UserID != uuid.Nil {
		query = query.Where("products.user_id = ?", *filter.UserID)
	}

	// Search by name or description
	if strings.TrimSpace(filter.Search) != "" {
		searchQuery := "%" + strings.ToLower(strings.TrimSpace(filter.Search)) + "%"
		query = query.Where("LOWER(products.name) LIKE ? OR LOWER(products.description) LIKE ?", searchQuery, searchQuery)
	}

	// Count total items before pagination
	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	switch filter.Sort {
	case "cheapest":
		query = query.Order("products.price ASC, products.created_at DESC")
	case "priciest":
		query = query.Order("products.price DESC, products.created_at DESC")
	case "latest":
		fallthrough
	default:
		query = query.Order("products.created_at DESC")
	}

	// Pagination
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var products []model.Product
	err := query.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) GetProductStats(productIDs []uuid.UUID, currentUserID *uuid.UUID) (likesCount map[uuid.UUID]int64, commentsCount map[uuid.UUID]int64, isLiked map[uuid.UUID]bool, err error) {
	likesCount = make(map[uuid.UUID]int64)
	commentsCount = make(map[uuid.UUID]int64)
	isLiked = make(map[uuid.UUID]bool)

	if len(productIDs) == 0 {
		return likesCount, commentsCount, isLiked, nil
	}

	// 1. Hitung jumlah likes per produk
	type CountResult struct {
		ProductID uuid.UUID `gorm:"column:product_id"`
		Total     int64     `gorm:"column:total"`
	}

	var likeCounts []CountResult
	err = r.db.Model(&model.Like{}).
		Select("product_id, COUNT(*) as total").
		Where("product_id IN ?", productIDs).
		Group("product_id").
		Scan(&likeCounts).Error
	if err != nil {
		return nil, nil, nil, err
	}
	for _, lc := range likeCounts {
		likesCount[lc.ProductID] = lc.Total
	}

	// 2. Hitung jumlah komentar per produk
	var commentCounts []CountResult
	err = r.db.Model(&model.Comment{}).
		Select("product_id, COUNT(*) as total").
		Where("product_id IN ?", productIDs).
		Group("product_id").
		Scan(&commentCounts).Error
	if err != nil {
		return nil, nil, nil, err
	}
	for _, cc := range commentCounts {
		commentsCount[cc.ProductID] = cc.Total
	}

	// 3. Cek apakah user sedang login telah memberi like
	if currentUserID != nil && *currentUserID != uuid.Nil {
		var likedProductIDs []uuid.UUID
		err = r.db.Model(&model.Like{}).
			Where("user_id = ? AND product_id IN ?", *currentUserID, productIDs).
			Pluck("product_id", &likedProductIDs).Error
		if err != nil {
			return nil, nil, nil, err
		}
		for _, pid := range likedProductIDs {
			isLiked[pid] = true
		}
	}

	return likesCount, commentsCount, isLiked, nil
}
