package service

import (
	"backend/internal/repository"
	"backend/pkg/apperror"
)

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryService interface {
	GetCategories() ([]CategoryResponse, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) GetCategories() ([]CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, apperror.Internal("Gagal mengambil data kategori: " + err.Error())
	}

	result := make([]CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
			Slug: c.Slug,
		})
	}

	return result, nil
}
