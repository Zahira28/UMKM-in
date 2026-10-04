package service_test

import (
	"testing"
	"time"

	"backend/internal/model"
	"backend/internal/service"

	"github.com/google/uuid"
)

type mockLikeRepo struct {
	likes map[string]*model.Like
}

func newMockLikeRepo() *mockLikeRepo {
	return &mockLikeRepo{
		likes: make(map[string]*model.Like),
	}
}

func (m *mockLikeRepo) key(userID, productID uuid.UUID) string {
	return userID.String() + ":" + productID.String()
}

func (m *mockLikeRepo) FindByUserAndProduct(userID uuid.UUID, productID uuid.UUID) (*model.Like, error) {
	like, exists := m.likes[m.key(userID, productID)]
	if !exists {
		return nil, nil
	}
	return like, nil
}

func (m *mockLikeRepo) Create(like *model.Like) error {
	like.CreatedAt = time.Now()
	m.likes[m.key(like.UserID, like.ProductID)] = like
	return nil
}

func (m *mockLikeRepo) Delete(userID uuid.UUID, productID uuid.UUID) error {
	delete(m.likes, m.key(userID, productID))
	return nil
}

func (m *mockLikeRepo) CountByProductID(productID uuid.UUID) (int64, error) {
	var count int64
	for _, l := range m.likes {
		if l.ProductID == productID {
			count++
		}
	}
	return count, nil
}

func TestLikeService_ToggleLike(t *testing.T) {
	mockProdRepo := newMockProductRepo()
	mockLikeRepo := newMockLikeRepo()
	likeService := service.NewLikeService(mockLikeRepo, mockProdRepo)

	productID := uuid.New()
	mockProdRepo.products[productID] = &model.Product{
		ID:   productID,
		Name: "Kue Nastar",
	}

	userID := uuid.New()

	t.Run("first toggle likes the product", func(t *testing.T) {
		res, err := likeService.ToggleLike(productID, userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !res.IsLiked {
			t.Errorf("expected is_liked true, got false")
		}
		if res.LikesCount != 1 {
			t.Errorf("expected likes_count 1, got %d", res.LikesCount)
		}
	})

	t.Run("second toggle unlikes the product", func(t *testing.T) {
		res, err := likeService.ToggleLike(productID, userID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.IsLiked {
			t.Errorf("expected is_liked false, got true")
		}
		if res.LikesCount != 0 {
			t.Errorf("expected likes_count 0, got %d", res.LikesCount)
		}
	})

	t.Run("fail toggle non-existent product", func(t *testing.T) {
		_, err := likeService.ToggleLike(uuid.New(), userID)
		if err == nil {
			t.Fatal("expected error for non-existent product, got nil")
		}
	})
}
