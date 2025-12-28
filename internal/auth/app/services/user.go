package services

import (
	"context"
	"fmt"
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	return s.userRepo.CreateUser(ctx, email, username, string(passwordHash))
}

func (s *UserService) LoginUser(ctx context.Context, user dto.RequestLoginUser) (dto.LoginUserDto, error) {
	userId, hashedPassword, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return dto.LoginUserDto{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return dto.LoginUserDto{}, fmt.Errorf("invalid email or password: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken()
	if err != nil {
		return dto.LoginUserDto{}, err
	}
	accessToken, err := s.jwtManager.GenerateAccessToken(userId, user.Email)
	if err != nil {
		return dto.LoginUserDto{}, err
	}

	expiresAt := time.Now().Add(s.jwtManager.RefreshTokenTTL())

	_, err = s.userRepo.CreateRefreshToken(ctx, userId, refreshToken, expiresAt)
	if err != nil {
		return dto.LoginUserDto{}, err
	}
	log.Print("user_id: ", userId)
	return dto.LoginUserDto{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       userId,
	}, nil
}

func (s *UserService) ValidateAccessToken(ctx context.Context, token string) (*jwt.Claims, error) {

	claims, err := s.jwtManager.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (string, error) {
	// Проверяем refresh token в БД
	userID, expiresAt, err := s.userRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Проверяем срок действия
	if time.Now().After(expiresAt) {
		// Удаляем истекший токен
		_ = s.userRepo.DeleteRefreshToken(ctx, refreshToken)
		return "", fmt.Errorf("refresh token expired")
	}

	// Получаем email пользователя
	email, err := s.userRepo.GetUserEmailByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Генерируем новый access token
	accessToken, err := s.jwtManager.GenerateAccessToken(userID, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) error {
	err := s.userRepo.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}
