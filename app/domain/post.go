package domain

import (
	"context"
	"time"
)

type Post struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	AuthorID  int64      `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type PostRepository interface {
	Create(ctx context.Context, tx Transaction, post *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	SelectForUpdate(ctx context.Context, tx Transaction, id int64) (*Post, error)
	Update(ctx context.Context, tx Transaction, id int64, post *Post) error
	GetAll(ctx context.Context, search SearchParam) ([]Post, int64, error)
}

type CreatePostRequestDTO struct {
	AuthorID int64  `json:"-"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type UpdatePostRequestDTO struct {
	AuthorID int64  `json:"-"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type CreatePostResponseDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  int64     `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdatePostResponseDTO struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	AuthorID  int64      `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DeletePostRequestDTO struct {
	AuthorID int64 `json:"-"`
}

type PostUseCase interface {
	Create(ctx context.Context, post *CreatePostRequestDTO) (*CreatePostResponseDTO, error)
	GetByID(ctx context.Context, id int64) (*Post, error)
	Update(ctx context.Context, id int64, post *UpdatePostRequestDTO) (*UpdatePostResponseDTO, error)
	Delete(ctx context.Context, id int64, post *DeletePostRequestDTO) error
	GetAll(ctx context.Context, search SearchParam) ([]Post, int64, error)
}
