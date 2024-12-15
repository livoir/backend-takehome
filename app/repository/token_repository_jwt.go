package repository

import (
	"app/domain"
	"app/pkg/common"
	"app/pkg/logger"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type TokenRepositoryJWT struct {
	privateKeyPath string
	publicKeyPath  string
}

func NewTokenRepositoryJWT(privateKey, publicKey string) domain.TokenRepository {
	return &TokenRepositoryJWT{
		privateKeyPath: privateKey,
		publicKeyPath:  publicKey,
	}
}

// Create implements domain.TokenRepository.
func (repository *TokenRepositoryJWT) Create(ctx context.Context, token *domain.TokenRequest) (string, error) {
	now := time.Now()
	claims := &jwt.MapClaims{
		"dat": token.Data,
		"exp": now.Add(token.ExpiresIn).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(repository.privateKeyPath))
	if err != nil {
		logger.Log.Error("failed to parse private key", zap.Error(err))
		return "", common.ErrInternalServerError
	}
	res, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		logger.Log.Error("failed to sign token", zap.Error(err))
		return "", common.ErrInternalServerError
	}
	return res, nil
}
