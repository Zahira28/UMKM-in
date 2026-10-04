package service_test

import (
	"errors"
	"testing"
	"time"

	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"backend/pkg/apperror"

	"github.com/google/uuid"
)

type mockProductRepo struct {
	products map[uuid.UUID]*model.Product
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{
		products: make(map[uuid.UUID]*model.Product),
	}
}

func (m *mockProductRepo) Create(product *model.Product) error {
	if product.ID == uuid.Nil {
		product.ID = uuid.New()
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepo) FindByID(id uuid.UUID) (*model.Product, error) {
	p, exists := m.products[id]
	if !exists {
		return nil, nil
	}
	return p, nil
}

func (m *mockProductRepo) FindAll(filter repository.ProductFilter) ([]model.Product, int64, error) {
	var list []model.Product
	for _, p := range m.products {
		list = append(list, *p)
	}
	return list, int64(len(list)), nil
}

func (m *mockProductRepo) Update(product *model.Product) error {
	m.products[product.ID] = product
	return nil
}

func (m *mockProductRepo) Delete(id uuid.UUID) error {
	delete(m.products, id)
	return nil
}

func (m *mockProductRepo) GetProductStats(productIDs []uuid.UUID, currentUserID *uuid.UUID) (map[uuid.UUID]int64, map[uuid.UUID]int64, map[uuid.UUID]bool, error) {
	likesCount := make(map[uuid.UUID]int64)
	commentsCount := make(map[uuid.UUID]int64)
	isLiked := make(map[uuid.UUID]bool)
	return likesCount, commentsCount, isLiked, nil
}

func TestProductService_CreateProduct(t *testing.T) {
	mockCatRepo := &mockCategoryRepo{
		categories: []model.Category{
			{ID: 1, Name: "Makanan", Slug: "makanan"},
		},
	}
	mockProdRepo := newMockProductRepo()
	prodService := service.NewProductService(mockProdRepo, mockCatRepo)

	userID := uuid.New()

	t.Run("success create product", func(t *testing.T) {
		req := service.CreateProductRequest{
			CategoryID:  1,
			Name:        "Kue Aren Lezat",
			Description: "Deskripsi kue aren",
			Price:       25000,
			Unit:        "box",
		}

		res, err := prodService.CreateProduct(userID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res == nil || res.ID == uuid.Nil {
			t.Fatalf("expected product ID to be generated")
		}
		if res.Name != "Kue Aren Lezat" {
			t.Errorf("expected name 'Kue Aren Lezat', got '%s'", res.Name)
		}
	})

	t.Run("fail name too short", func(t *testing.T) {
		req := service.CreateProductRequest{
			CategoryID: 1,
			Name:       "Ku",
			Price:      25000,
			Unit:       "box",
		}

		_, err := prodService.CreateProduct(userID, req)
		if err == nil {
			t.Fatal("expected error for short name, got nil")
		}
	})

	t.Run("fail negative price", func(t *testing.T) {
		req := service.CreateProductRequest{
			CategoryID: 1,
			Name:       "Kue Aren",
			Price:      -1000,
			Unit:       "box",
		}

		_, err := prodService.CreateProduct(userID, req)
		if err == nil {
			t.Fatal("expected error for negative price, got nil")
		}
	})

	t.Run("fail invalid category", func(t *testing.T) {
		req := service.CreateProductRequest{
			CategoryID: 999,
			Name:       "Kue Aren",
			Price:      25000,
			Unit:       "box",
		}

		_, err := prodService.CreateProduct(userID, req)
		if err == nil {
			t.Fatal("expected error for invalid category, got nil")
		}
	})
}

func TestProductService_Ownership(t *testing.T) {
	mockCatRepo := &mockCategoryRepo{
		categories: []model.Category{
			{ID: 1, Name: "Makanan", Slug: "makanan"},
		},
	}
	mockProdRepo := newMockProductRepo()
	prodService := service.NewProductService(mockProdRepo, mockCatRepo)

	ownerID := uuid.New()
	otherUserID := uuid.New()

	productID := uuid.New()
	mockProdRepo.products[productID] = &model.Product{
		ID:          productID,
		UserID:      ownerID,
		CategoryID:  1,
		Name:        "Kue Aren Owner",
		Price:       25000,
		Unit:        "box",
		IsAvailable: true,
	}

	t.Run("forbidden update if not owner", func(t *testing.T) {
		newName := "Kue Aren Bajakan"
		req := service.UpdateProductRequest{
			Name: &newName,
		}

		_, err := prodService.UpdateProduct(productID, otherUserID, req)
		if err == nil {
			t.Fatal("expected forbidden error, got nil")
		}

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.StatusCode != 403 {
				t.Errorf("expected status 403, got %d", appErr.StatusCode)
			}
		} else {
			t.Errorf("expected AppError, got %T", err)
		}
	})

	t.Run("forbidden delete if not owner", func(t *testing.T) {
		err := prodService.DeleteProduct(productID, otherUserID)
		if err == nil {
			t.Fatal("expected forbidden error, got nil")
		}

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.StatusCode != 403 {
				t.Errorf("expected status 403, got %d", appErr.StatusCode)
			}
		} else {
			t.Errorf("expected AppError, got %T", err)
		}
	})

	t.Run("successful update by owner", func(t *testing.T) {
		newName := "Kue Aren Spesial Owner"
		req := service.UpdateProductRequest{
			Name: &newName,
		}

		updated, err := prodService.UpdateProduct(productID, ownerID, req)
		if err != nil {
			t.Fatalf("expected successful update, got %v", err)
		}
		if updated.Name != newName {
			t.Errorf("expected name '%s', got '%s'", newName, updated.Name)
		}
	})

	t.Run("successful delete by owner", func(t *testing.T) {
		err := prodService.DeleteProduct(productID, ownerID)
		if err != nil {
			t.Fatalf("expected successful delete, got %v", err)
		}
		if _, exists := mockProdRepo.products[productID]; exists {
			t.Fatal("expected product to be deleted from repo")
		}
	})
}
