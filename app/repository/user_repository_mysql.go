package repository

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"
	"context"
	"database/sql"

	"go.uber.org/zap"
)

type UserRepositoryMySQL struct {
	sql *sql.DB
}

func NewUserRepositoryMySQL(sql *sql.DB) domain.UserRepository {
	return &UserRepositoryMySQL{
		sql: sql,
	}
}

// Create implements domain.UserRepository.
func (repository *UserRepositoryMySQL) Create(ctx context.Context, tx domain.Transaction, user *domain.User) error {
	result, err := tx.GetTx().ExecContext(ctx, "INSERT INTO users (email, name, password_hash) VALUES (?, ?, ?)", user.Email, user.Name, user.PasswordHash)
	if err != nil {
		logger.Log.Error("failed to insert user", zap.Error(err))
		return common.ErrInternalServerError
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Log.Error("failed to get last insert id", zap.Error(err))
		return common.ErrInternalServerError
	}
	user.ID = id
	return nil
}

// FindByEmail implements domain.UserRepository.
func (repository *UserRepositoryMySQL) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := repository.sql.QueryRowContext(ctx, "SELECT id, name, email, password_hash, created_at, updated_at, deleted_at FROM users WHERE email = ? and deleted_at is NULL", email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrEmailNotFound
		}
		logger.Log.Error("failed to select user by email", zap.Error(err))
		return nil, common.ErrInternalServerError
	}
	return &user, nil
}

// FindByID implements domain.UserRepository.
func (repository *UserRepositoryMySQL) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	err := repository.sql.QueryRowContext(ctx, "SELECT id, name, email, password_hash FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrUserNotFound
		}
		logger.Log.Error("failed to select user by id", zap.Error(err))
		return nil, common.ErrInternalServerError
	}
	return &user, nil
}
