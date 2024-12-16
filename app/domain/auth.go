package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int64      `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

type RegisterRequestDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponseDTO struct {
	ID int64 `json:"id"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	*User
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
}

type TokenRequest struct {
	Data      interface{}   `json:"data"`
	ExpiresIn time.Duration `json:"expires_in"`
}

type VerifyTokenResponse struct {
	UserID int64 `json:"user_id"`
}

type AuthUseCase interface {
	Register(ctx context.Context, request *RegisterRequestDTO) (*RegisterResponseDTO, error)
	Login(ctx context.Context, request *LoginRequestDTO) (*LoginResponseDTO, error)
	VerifyToken(ctx context.Context, token string) (*VerifyTokenResponse, error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
	Create(ctx context.Context, tx Transaction, user *User) error
}

type TokenRepository interface {
	Create(ctx context.Context, token *TokenRequest) (string, error)
	Verify(ctx context.Context, token string) (*VerifyTokenResponse, error)
}
