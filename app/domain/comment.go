package domain

import (
	"context"
	"time"
)

type Comment struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	PostID     int64     `json:"post_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type CommentRepository interface {
	Create(ctx context.Context, tx Transaction, comment *Comment) error
	FindByPostID(ctx context.Context, postID int64, param SearchParam) ([]*Comment, int64, error)
}

type CreateCommentRequestDTO struct {
	AuthorID int64  `json:"-"`
	Content  string `json:"content" binding:"required"`
}

type CreateCommentResponseDTO struct {
	ID         int64     `json:"id"`
	Content    string    `json:"content"`
	PostID     int64     `json:"post_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type CommentUsecase interface {
	CreateComment(ctx context.Context, postId int64, req CreateCommentRequestDTO) (*CreateCommentResponseDTO, error)
	FindCommentsByPostID(ctx context.Context, postID int64, param SearchParam) ([]*Comment, int64, error)
}
