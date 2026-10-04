package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperror"

	"github.com/google/uuid"
)

type ToggleLikeResponse struct {
	IsLiked    bool  `json:"is_liked"`
	LikesCount int64 `json:"likes_count"`
}

type LikeService interface {
	ToggleLike(productID uuid.UUID, userID uuid.UUID) (*ToggleLikeResponse, error)
}

type likeService struct {
	likeRepo    repository.LikeRepository
	productRepo repository.ProductRepository
}

func NewLikeService(likeRepo repository.LikeRepository, productRepo repository.ProductRepository) LikeService {
	return &likeService{
		likeRepo:    likeRepo,
		productRepo: productRepo,
	}
}

func (s *likeService) ToggleLike(productID uuid.UUID, userID uuid.UUID) (*ToggleLikeResponse, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, apperror.Internal("Gagal memuat produk: " + err.Error())
	}
	if product == nil {
		return nil, apperror.NotFound("Produk tidak ditemukan")
	}

	existingLike, err := s.likeRepo.FindByUserAndProduct(userID, productID)
	if err != nil {
		return nil, apperror.Internal("Gagal memeriksa status like: " + err.Error())
	}

	isLiked := false
	if existingLike != nil {
		// Unlike
		if err := s.likeRepo.Delete(userID, productID); err != nil {
			return nil, apperror.Internal("Gagal membatalkan like: " + err.Error())
		}
		isLiked = false
	} else {
		// Like
		newLike := &model.Like{
			ProductID: productID,
			UserID:    userID,
		}
		if err := s.likeRepo.Create(newLike); err != nil {
			return nil, apperror.Internal("Gagal menambahkan like: " + err.Error())
		}
		isLiked = true
	}

	count, err := s.likeRepo.CountByProductID(productID)
	if err != nil {
		return nil, apperror.Internal("Gagal menghitung jumlah like: " + err.Error())
	}

	return &ToggleLikeResponse{
		IsLiked:    isLiked,
		LikesCount: count,
	}, nil
}
