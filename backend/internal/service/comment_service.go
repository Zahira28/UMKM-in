package service

import (
	"strings"
	"time"

	"backend/internal/model"
	"backend/internal/repository"
	"backend/pkg/apperror"

	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	Content  string  `json:"content"`
	ParentID *uint64 `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type CommentUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
}

type CommentReplyResponse struct {
	ID        uint64               `json:"id"`
	ParentID  *uint64              `json:"parent_id"`
	Content   string               `json:"content"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	User      *CommentUserResponse `json:"user"`
}

type CommentResponse struct {
	ID        uint64                 `json:"id"`
	Content   string                 `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	User      *CommentUserResponse   `json:"user"`
	Replies   []CommentReplyResponse `json:"replies"`
}

type CommentCreatedResponse struct {
	ID        uint64    `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Content   string    `json:"content"`
	ParentID  *uint64   `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentService interface {
	GetComments(productID uuid.UUID) ([]CommentResponse, error)
	CreateComment(productID uuid.UUID, userID uuid.UUID, req CreateCommentRequest) (*CommentCreatedResponse, error)
	UpdateComment(commentID uint64, userID uuid.UUID, req UpdateCommentRequest) (*CommentReplyResponse, error)
	DeleteComment(commentID uint64, userID uuid.UUID) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	productRepo repository.ProductRepository
}

func NewCommentService(commentRepo repository.CommentRepository, productRepo repository.ProductRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		productRepo: productRepo,
	}
}

func (s *commentService) GetComments(productID uuid.UUID) ([]CommentResponse, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, apperror.Internal("Gagal memuat produk: " + err.Error())
	}
	if product == nil {
		return nil, apperror.NotFound("Produk tidak ditemukan")
	}

	comments, err := s.commentRepo.FindByProductID(productID)
	if err != nil {
		return nil, apperror.Internal("Gagal mengambil komentar produk: " + err.Error())
	}

	result := make([]CommentResponse, 0, len(comments))
	for _, c := range comments {
		var userResp *CommentUserResponse
		if c.User != nil {
			userResp = &CommentUserResponse{
				ID:        c.User.ID,
				Username:  c.User.Username,
				FullName:  c.User.FullName,
				AvatarURL: c.User.AvatarURL,
			}
		}

		repliesResp := make([]CommentReplyResponse, 0, len(c.Replies))
		for _, r := range c.Replies {
			var replyUserResp *CommentUserResponse
			if r.User != nil {
				replyUserResp = &CommentUserResponse{
					ID:        r.User.ID,
					Username:  r.User.Username,
					FullName:  r.User.FullName,
					AvatarURL: r.User.AvatarURL,
				}
			}

			repliesResp = append(repliesResp, CommentReplyResponse{
				ID:        r.ID,
				ParentID:  r.ParentID,
				Content:   r.Content,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
				User:      replyUserResp,
			})
		}

		result = append(result, CommentResponse{
			ID:        c.ID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			User:      userResp,
			Replies:   repliesResp,
		})
	}

	return result, nil
}

func (s *commentService) CreateComment(productID uuid.UUID, userID uuid.UUID, req CreateCommentRequest) (*CommentCreatedResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, apperror.BadRequest("Isi komentar tidak boleh kosong")
	}

	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, apperror.Internal("Gagal memuat produk: " + err.Error())
	}
	if product == nil {
		return nil, apperror.NotFound("Produk tidak ditemukan")
	}

	var parentID *uint64
	if req.ParentID != nil && *req.ParentID > 0 {
		parentComment, err := s.commentRepo.FindByID(*req.ParentID)
		if err != nil {
			return nil, apperror.Internal("Gagal memvalidasi komentar utama: " + err.Error())
		}
		if parentComment == nil {
			return nil, apperror.NotFound("Komentar utama yang ingin dibalas tidak ditemukan")
		}
		if parentComment.ProductID != productID {
			return nil, apperror.BadRequest("Komentar tidak berada pada produk yang sama")
		}

		// Jika komentar yang dibalas adalah balasan, jadikan parent ID dari komentar utamanya (flat 1-level)
		if parentComment.ParentID != nil {
			parentID = parentComment.ParentID
		} else {
			targetID := parentComment.ID
			parentID = &targetID
		}
	}

	comment := &model.Comment{
		ProductID: productID,
		UserID:    userID,
		ParentID:  parentID,
		Content:   content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, apperror.Internal("Gagal menyimpan komentar: " + err.Error())
	}

	return &CommentCreatedResponse{
		ID:        comment.ID,
		ProductID: comment.ProductID,
		Content:   comment.Content,
		ParentID:  comment.ParentID,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (s *commentService) UpdateComment(commentID uint64, userID uuid.UUID, req UpdateCommentRequest) (*CommentReplyResponse, error) {
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, apperror.BadRequest("Isi komentar tidak boleh kosong")
	}

	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return nil, apperror.Internal("Gagal memuat komentar: " + err.Error())
	}
	if comment == nil {
		return nil, apperror.NotFound("Komentar tidak ditemukan")
	}

	// Ownership check: hanya pembuat komentar yang dapat mengedit
	if comment.UserID != userID {
		return nil, apperror.Forbidden("Anda tidak memiliki akses untuk mengubah komentar ini")
	}

	comment.Content = content
	if err := s.commentRepo.Update(comment); err != nil {
		return nil, apperror.Internal("Gagal memperbarui komentar: " + err.Error())
	}

	var userResp *CommentUserResponse
	if comment.User != nil {
		userResp = &CommentUserResponse{
			ID:        comment.User.ID,
			Username:  comment.User.Username,
			FullName:  comment.User.FullName,
			AvatarURL: comment.User.AvatarURL,
		}
	}

	return &CommentReplyResponse{
		ID:        comment.ID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User:      userResp,
	}, nil
}

func (s *commentService) DeleteComment(commentID uint64, userID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return apperror.Internal("Gagal memuat komentar: " + err.Error())
	}
	if comment == nil {
		return apperror.NotFound("Komentar tidak ditemukan")
	}

	// Hak akses: pembuat komentar ATAU pemilik produk yang dikomentari
	isAuthor := comment.UserID == userID
	isProductOwner := comment.Product != nil && comment.Product.UserID == userID

	if !isAuthor && !isProductOwner {
		return apperror.Forbidden("Anda tidak memiliki izin untuk menghapus komentar ini")
	}

	if err := s.commentRepo.Delete(commentID); err != nil {
		return apperror.Internal("Gagal menghapus komentar: " + err.Error())
	}

	return nil
}
