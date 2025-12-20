package services

import (
	"context"
	"habr/internal/auth/app/repositories"
	"habr/internal/auth/core/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*repositories.UserRepository
	*jwt.JWTManager
}

func NewUserService(repo *repositories.UserRepository, jwt *jwt.JWTManager) *UserService {
	return &UserService{UserRepository: repo, JWTManager: jwt}
}

func (s *UserService) RegisterUser(ctx context.Context, email, username, password string) (int64, error) {
	// Хэширование пароля
	passwordHash, err := hashPassword(password)
	if err != nil {
		return -1, err
	}

	tokenHash, err := s.JWTManager.GenerateRefreshToken()
	if err != nil {
		return -1, err
	}

	return s.UserRepository.CreateUser(ctx, email, username, passwordHash, tokenHash)
}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (int64, error) {
	// Хэширование пароля
	passwordHash, err := hashPassword(password)
	if err != nil {
		return -1, err
	}

	return s.UserRepository.LoginUser(ctx, email, passwordHash)
}

func hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MaxCost)
	if err != nil {
		return "", err
	}

	return string(passwordHash), nil
}
