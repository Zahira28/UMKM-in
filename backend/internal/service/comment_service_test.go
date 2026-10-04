package service_test

import (
	"errors"
	"testing"
	"time"

	"backend/internal/model"
	"backend/internal/service"
	"backend/pkg/apperror"

	"github.com/google/uuid"
)

type mockCommentRepo struct {
	comments map[uint64]*model.Comment
	nextID   uint64
}

func newMockCommentRepo() *mockCommentRepo {
	return &mockCommentRepo{
		comments: make(map[uint64]*model.Comment),
		nextID:   1,
	}
}

func (m *mockCommentRepo) Create(comment *model.Comment) error {
	comment.ID = m.nextID
	m.nextID++
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepo) FindByID(id uint64) (*model.Comment, error) {
	c, exists := m.comments[id]
	if !exists {
		return nil, nil
	}
	return c, nil
}

func (m *mockCommentRepo) FindByProductID(productID uuid.UUID) ([]model.Comment, error) {
	var list []model.Comment
	for _, c := range m.comments {
		if c.ProductID == productID && c.ParentID == nil {
			topComment := *c
			for _, reply := range m.comments {
				if reply.ParentID != nil && *reply.ParentID == c.ID {
					topComment.Replies = append(topComment.Replies, *reply)
				}
			}
			list = append(list, topComment)
		}
	}
	return list, nil
}

func (m *mockCommentRepo) Update(comment *model.Comment) error {
	comment.UpdatedAt = time.Now()
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepo) Delete(id uint64) error {
	delete(m.comments, id)
	return nil
}

func TestCommentService_CreateComment(t *testing.T) {
	mockProdRepo := newMockProductRepo()
	mockComRepo := newMockCommentRepo()
	commentService := service.NewCommentService(mockComRepo, mockProdRepo)

	sellerID := uuid.New()
	productID := uuid.New()
	mockProdRepo.products[productID] = &model.Product{
		ID:     productID,
		UserID: sellerID,
		Name:   "Kue Nastar",
	}

	userID := uuid.New()

	t.Run("success create top-level comment", func(t *testing.T) {
		req := service.CreateCommentRequest{
			Content: "Bisa kirim hari ini?",
		}

		res, err := commentService.CreateComment(productID, userID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.ID == 0 || res.Content != req.Content {
			t.Errorf("unexpected comment response: %+v", res)
		}
	})

	t.Run("success reply to existing comment", func(t *testing.T) {
		parentID := uint64(1)
		req := service.CreateCommentRequest{
			Content:  "Bisa kak, ready stock!",
			ParentID: &parentID,
		}

		res, err := commentService.CreateComment(productID, sellerID, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.ParentID == nil || *res.ParentID != parentID {
			t.Errorf("expected parent_id %d, got %v", parentID, res.ParentID)
		}
	})

	t.Run("fail empty content", func(t *testing.T) {
		req := service.CreateCommentRequest{
			Content: "   ",
		}

		_, err := commentService.CreateComment(productID, userID, req)
		if err == nil {
			t.Fatal("expected error for empty content, got nil")
		}
	})

	t.Run("fail product not found", func(t *testing.T) {
		req := service.CreateCommentRequest{
			Content: "Halo",
		}

		_, err := commentService.CreateComment(uuid.New(), userID, req)
		if err == nil {
			t.Fatal("expected not found error, got nil")
		}
	})
}

func TestCommentService_UpdateAndDelete(t *testing.T) {
	mockProdRepo := newMockProductRepo()
	mockComRepo := newMockCommentRepo()
	commentService := service.NewCommentService(mockComRepo, mockProdRepo)

	sellerID := uuid.New()
	authorID := uuid.New()
	otherUserID := uuid.New()

	product := &model.Product{
		ID:     uuid.New(),
		UserID: sellerID,
		Name:   "Kue Nastar",
	}
	mockProdRepo.products[product.ID] = product

	comment := &model.Comment{
		ID:        1,
		ProductID: product.ID,
		UserID:    authorID,
		Content:   "Apakah masih ada stok?",
		Product:   product,
	}
	mockComRepo.comments[comment.ID] = comment

	t.Run("forbidden update if not author", func(t *testing.T) {
		req := service.UpdateCommentRequest{Content: "Diedit orang lain"}
		_, err := commentService.UpdateComment(1, otherUserID, req)
		if err == nil {
			t.Fatal("expected forbidden error, got nil")
		}

		var appErr *apperror.AppError
		if errors.As(err, &appErr) && appErr.StatusCode != 403 {
			t.Errorf("expected status 403, got %d", appErr.StatusCode)
		}
	})

	t.Run("successful update by author", func(t *testing.T) {
		req := service.UpdateCommentRequest{Content: "Apakah masih ada stok rasa cokelat?"}
		res, err := commentService.UpdateComment(1, authorID, req)
		if err != nil {
			t.Fatalf("expected successful update, got %v", err)
		}
		if res.Content != req.Content {
			t.Errorf("expected content '%s', got '%s'", req.Content, res.Content)
		}
	})

	t.Run("forbidden delete by third-party user", func(t *testing.T) {
		err := commentService.DeleteComment(1, otherUserID)
		if err == nil {
			t.Fatal("expected forbidden error, got nil")
		}
	})

	t.Run("successful delete by product seller (store owner moderation)", func(t *testing.T) {
		err := commentService.DeleteComment(1, sellerID)
		if err != nil {
			t.Fatalf("expected seller to be able to delete comment on own product, got %v", err)
		}
		if _, exists := mockComRepo.comments[1]; exists {
			t.Fatal("expected comment 1 to be deleted")
		}
	})
}
