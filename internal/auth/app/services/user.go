package services

import (
	"context"
	"habr/internal/auth/app/repositories"
	"habr/internal/auth/core/jwt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo   *repositories.UserRepository
	jwtManager *jwt.Manager
}

func NewUserService(repo *repositories.UserRepository, jwt *jwt.Manager) *UserService {
	return &UserService{userRepo: repo, jwtManager: jwt}
}

func (s *UserService) RegisterUser(ctx context.Context, email, username, password string) (int64, error) {
	// Хэширование пароля
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return -1, err
	}

	tokenHash, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return -1, err
	}
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 дней
	return s.userRepo.CreateUser(ctx, email, username, string(passwordHash), tokenHash, expiresAt)
}

func (s *UserService) LoginUser(ctx context.Context, email string, password string) (int64, error) {
	// Хэширование пароля
	// passwordHash, err := hashPassword(password)

	userId, hashedPassword, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return -1, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return -1, err
	}
	log.Print("user_id", userId)
	return userId, err
}
