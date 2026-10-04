package repository

import (
	"errors"

	"backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeRepository interface {
	FindByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (*model.Like, error)
	Create(like *model.Like) error
	Delete(userID uuid.UUID, productID uuid.UUID) error
	CountByProductID(productID uuid.UUID) (int64, error)
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) FindByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (*model.Like, error) {
	var like model.Like
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepository) Create(like *model.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(userID uuid.UUID, productID uuid.UUID) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.Like{}).Error
}

func (r *likeRepository) CountByProductID(productID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Like{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}
