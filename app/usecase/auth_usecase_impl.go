package usecase

import (
	"app/domain"
	"app/pkg/common"
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseImpl struct {
	userRepository  domain.UserRepository
	tokenRepository domain.TokenRepository
	transactor      domain.Transactor
}

func NewAuthUseCaseImpl(userRepository domain.UserRepository, tokenRepository domain.TokenRepository, transactor domain.Transactor) domain.AuthUseCase {
	return &AuthUseCaseImpl{
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		transactor:      transactor,
	}
}

// Login implements domain.AuthUseCase.
func (uc *AuthUseCaseImpl) Login(ctx context.Context, request *domain.LoginRequestDTO) (*domain.LoginResponseDTO, error) {
	user, err := uc.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, common.ErrInvalidPassword
	}
	data := map[string]interface{}{
		"id": user.ID,
	}
	tokenRequest := &domain.TokenRequest{
		Data:      data,
		ExpiresIn: time.Duration(10) * time.Minute,
	}
	token, err := uc.tokenRepository.Create(ctx, tokenRequest)
	if err != nil {
		return nil, err
	}
	tokenRequest.ExpiresIn = time.Duration(24*7) * time.Hour
	refreshToken, err := uc.tokenRepository.Create(ctx, tokenRequest)
	if err != nil {
		return nil, err
	}
	response := &domain.LoginResponseDTO{
		User:         user,
		AccessToken:  token,
		RefreshToken: refreshToken,
	}
	return response, nil
}

// Register implements domain.AuthUseCase.
func (uc *AuthUseCaseImpl) Register(ctx context.Context, request *domain.RegisterRequestDTO) (*domain.RegisterResponseDTO, error) {
	now := time.Now()
	user, err := uc.userRepository.FindByEmail(ctx, request.Email)
	if err == common.ErrEmailNotFound {
		err = nil
	}
	if err != nil {

		return nil, err
	}
	if user != nil {
		return nil, common.ErrEmailAlreadyExists
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), 14)
	if err != nil {
		return nil, err
	}
	user = &domain.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
	}
	tx, err := uc.transactor.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	err = uc.userRepository.Create(ctx, tx, user)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	response := &domain.RegisterResponseDTO{ID: user.ID}
	return response, nil
}

// VerifyToken implements domain.AuthUseCase.
func (uc *AuthUseCaseImpl) VerifyToken(ctx context.Context, token string) (*domain.VerifyTokenResponse, error) {
	return uc.tokenRepository.Verify(ctx, token)
}

// RefreshToken implements domain.AuthUseCase.
func (uc *AuthUseCaseImpl) RefreshToken(ctx context.Context, token string) (*domain.RefreshTokenResponse, error) {
	verifiedToken, err := uc.tokenRepository.Verify(ctx, token)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"id": verifiedToken.UserID,
	}
	tokenRequest := &domain.TokenRequest{
		Data:      data,
		ExpiresIn: time.Duration(10) * time.Minute,
	}
	newToken, err := uc.tokenRepository.Create(ctx, tokenRequest)
	if err != nil {
		return nil, err
	}
	return &domain.RefreshTokenResponse{
		AccessToken: newToken,
	}, nil
}
