package repository

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type CommentRepositoryMySQL struct {
	db *sql.DB
}

func NewCommentRepositoryMySQL(db *sql.DB) domain.CommentRepository {
	return &CommentRepositoryMySQL{db: db}
}

// Create implements domain.CommentRepository.
func (repository *CommentRepositoryMySQL) Create(ctx context.Context, tx domain.Transaction, comment *domain.Comment) error {
	result, err := tx.GetTx().ExecContext(ctx, "INSERT INTO comments (content, post_id, author_name) VALUES (?, ?, ?)", comment.Content, comment.PostID, comment.AuthorName)
	if err != nil {
		logger.Log.Error("failed to insert comment", zap.Error(err))
		return common.ErrInternalServerError
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("failed to get last insert id", zap.Error(err))
		return common.ErrInternalServerError
	}
	comment.ID = id
	return nil
}

// FindByPostID implements domain.CommentRepository.
func (repository *CommentRepositoryMySQL) FindByPostID(ctx context.Context, postID int64, param domain.SearchParam) ([]*domain.Comment, int64, error) {
	var comments []*domain.Comment
	query := "SELECT count(id) FROM comments WHERE post_id = ?"
	row := repository.db.QueryRowContext(ctx, query, postID)
	var total int64
	if err := row.Scan(&total); err != nil {
		logger.Log.Error("failed to count comments", zap.Error(err))
		return nil, 0, common.ErrInternalServerError
	}
	query = "SELECT id, content, post_id, author_name, created_at FROM comments WHERE post_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := repository.db.QueryContext(ctx, query, postID, param.Limit, (param.Page-1)*param.Limit)
	if err != nil {
		logger.Log.Error("failed to select comments", zap.Error(err))
		return nil, 0, common.ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var comment domain.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.AuthorName, &comment.CreatedAt)
		if err != nil {
			logger.Log.Error("failed to scan comment", zap.Error(err))
			return nil, 0, common.ErrInternalServerError
		}
		comments = append(comments, &comment)
	}
	return comments, total, nil
}
