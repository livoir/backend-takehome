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
	privateKey string
	publicKey  string
}

func NewTokenRepositoryJWT(privateKey, publicKey string) domain.TokenRepository {
	return &TokenRepositoryJWT{
		privateKey: privateKey,
		publicKey:  publicKey,
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
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(repository.privateKey))
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

// Verify implements domain.TokenRepository.
func (repository *TokenRepositoryJWT) Verify(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(repository.publicKey))
	if err != nil {
		logger.Log.Error("failed to parse public key", zap.Error(err))
		return nil, common.ErrInternalServerError
	}
	tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, common.ErrInvalidTokenMethod
		}
		return key, nil
	})
	if err != nil {
		logger.Log.Error("failed to parse token", zap.Error(err))
		return nil, common.ErrInvalidToken
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, common.ErrInvalidToken
	}
	userID, ok := claims["dat"].(map[string]interface{})["id"].(float64)
	if !ok {
		return nil, common.ErrInvalidToken
	}
	return &domain.VerifyTokenResponse{
		UserID: int64(userID),
	}, nil
}
