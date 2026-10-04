package service_test

import (
	"testing"

	"backend/internal/model"
	"backend/internal/service"
)

type mockCategoryRepo struct {
	categories []model.Category
	err        error
}

func (m *mockCategoryRepo) FindAll() ([]model.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func (m *mockCategoryRepo) FindByID(id uint) (*model.Category, error) {
	for _, c := range m.categories {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryRepo) FindBySlug(slug string) (*model.Category, error) {
	for _, c := range m.categories {
		if c.Slug == slug {
			return &c, nil
		}
	}
	return nil, nil
}

func TestCategoryService_GetCategories(t *testing.T) {
	mockRepo := &mockCategoryRepo{
		categories: []model.Category{
			{ID: 1, Name: "Makanan", Slug: "makanan"},
			{ID: 2, Name: "Minuman", Slug: "minuman"},
			{ID: 3, Name: "Pakaian", Slug: "pakaian"},
		},
	}

	catService := service.NewCategoryService(mockRepo)
	categories, err := catService.GetCategories()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(categories) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(categories))
	}

	if categories[0].Name != "Makanan" || categories[0].Slug != "makanan" {
		t.Errorf("unexpected category at index 0: %+v", categories[0])
	}
}
