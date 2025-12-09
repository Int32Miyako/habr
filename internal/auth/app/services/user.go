package services

import (
	"context"
	"habr/internal/auth/app/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, email, username, passwordHash string) (int64, error) {
	return s.userRepo.CreateUser(ctx, username, passwordHash, email)
}
