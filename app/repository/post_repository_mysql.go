package repository

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

type PostRepositoryMySQL struct {
	db *sql.DB
}

func NewPostRepositoryMySQL(db *sql.DB) domain.PostRepository {
	return &PostRepositoryMySQL{db: db}
}

// Create implements domain.PostRepository.
func (repository *PostRepositoryMySQL) Create(ctx context.Context, tx domain.Transaction, post *domain.Post) error {
	result, err := tx.GetTx().ExecContext(ctx, "INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)", post.Title, post.Content, post.AuthorID)
	if err != nil {
		logger.Log.Error("failed to insert user", zap.Error(err))
		return common.ErrInternalServerError
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("failed to get last insert id", zap.Error(err))
		return common.ErrInternalServerError
	}
	post.ID = id
	return nil
}

// GetAll implements domain.PostRepository.
func (repository *PostRepositoryMySQL) GetAll(ctx context.Context, search domain.SearchParam) ([]domain.Post, error) {
	var posts []domain.Post
	query := "SELECT id, title, content, author_id, created_at, updated_at, deleted_at FROM posts WHERE %s LIMIT ? OFFSET ?"
	if search.Search != "" {
		query = fmt.Sprintf(query, "title LIKE ? or content LIKE ?")
	}
	rows, err := repository.db.QueryContext(ctx, query, search.Limit, search.Limit*(search.Page-1))
	if err != nil {
		logger.Log.Error("failed to query posts", zap.Error(err))
		return nil, common.ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt); err != nil {
			logger.Log.Error("failed to scan post", zap.Error(err))
			return nil, common.ErrInternalServerError
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetByID implements domain.PostRepository.
func (repository *PostRepositoryMySQL) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	var post domain.Post
	err := repository.db.QueryRowContext(ctx, "SELECT id, title, content, author_id, created_at, updated_at, deleted_at FROM posts WHERE id = ?", id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrPostNotFound
		}
		logger.Log.Error("failed to select post by id", zap.Error(err))
		return nil, common.ErrInternalServerError
	}
	return &post, nil
}

// SelectForUpdate implements domain.PostRepository.
func (repository *PostRepositoryMySQL) SelectForUpdate(ctx context.Context, tx domain.Transaction, id int64) (*domain.Post, error) {
	var post domain.Post
	err := tx.GetTx().QueryRowContext(ctx, "SELECT id, title, content, author_id, created_at, updated_at, deleted_at FROM posts WHERE id = ? FOR UPDATE", id).Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrPostNotFound
		}
		logger.Log.Error("failed to select post for update", zap.Error(err))
		return nil, common.ErrInternalServerError
	}
	return &post, nil
}

// Update implements domain.PostRepository.
func (repository *PostRepositoryMySQL) Update(ctx context.Context, tx domain.Transaction, id int64, post *domain.Post) error {
	result, err := tx.GetTx().ExecContext(ctx, "UPDATE posts SET title = ?, content = ?, updated_at = ?, deleted_at = ? WHERE id = ?", post.Title, post.Content, post.DeletedAt, id)
	if err != nil {
		logger.Log.Error("failed to update post", zap.Error(err))
		return common.ErrInternalServerError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Log.Error("failed to get rows affected", zap.Error(err))
		return common.ErrInternalServerError
	}
	if rowsAffected == 0 {
		return common.ErrPostNotFound
	}
	return nil
}