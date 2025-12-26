package services

import (
	"context"
	"habr/internal/auth/app/repositories"
	"habr/internal/auth/core/jwt"
	"habr/internal/blog/http-server/dto"
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

	return s.userRepo.CreateUser(ctx, email, username, string(passwordHash))
}

func (s *UserService) LoginUser(ctx context.Context, user dto.RequestLoginUser) (dto.ResponseLoginUser, error) {
	userId, hashedPassword, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))

	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}
	accessToken, err := s.jwtManager.GenerateAccessToken(userId, user.Email)
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}

	expiresAt := time.Now().Add(s.jwtManager.RefreshTokenTTL() * 24 * time.Hour) // 30 дней

	_, err = s.userRepo.CreateRefreshToken(ctx, userId, refreshToken, expiresAt)
	if err != nil {
		return dto.ResponseLoginUser{}, err
	}
	log.Print("user_id", userId)
	return dto.ResponseLoginUser{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       userId,
	}, err
}
