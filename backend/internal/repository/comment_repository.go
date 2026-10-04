package repository

import (
	"errors"

	"backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *model.Comment) error
	FindByID(id uint64) (*model.Comment, error)
	FindByProductID(productID uuid.UUID) ([]model.Comment, error)
	Update(comment *model.Comment) error
	Delete(id uint64) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint64) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").Preload("Product").First(&comment, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByProductID(productID uuid.UUID) ([]model.Comment, error) {
	var comments []model.Comment
	// Ambil komentar utama (parent_id IS NULL) beserta replies dan relasi user masing-masing
	err := r.db.
		Where("product_id = ? AND parent_id IS NULL", productID).
		Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC").Preload("User")
		}).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Update(comment *model.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Comment{}, "id = ?", id).Error
}
