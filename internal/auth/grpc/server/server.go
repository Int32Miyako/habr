package server

import (
	"context"
	"fmt"
	"habr/internal/auth/app/services"
	"habr/protos/gen/go/auth"
	"log/slog"

	"google.golang.org/grpc"
)

type serverAPI struct {
	auth.UnimplementedAuthServer
	userService *services.UserService
	log         *slog.Logger
}

func Register(gRPCServer *grpc.Server, userService *services.UserService, log *slog.Logger) {
	auth.RegisterAuthServer(gRPCServer, &serverAPI{
		userService: userService,
		log:         log,
	})
}

// Register реализует метод регистрации нового пользователя
func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	const op = "grpc.Register"

	// Валидация входных данных
	if req.GetEmail() == "" {
		return nil, fmt.Errorf("%s: email is required", op)
	}
	if req.GetUsername() == "" {
		return nil, fmt.Errorf("%s: username is required", op)
	}
	if req.GetPassword() == "" {
		return nil, fmt.Errorf("%s: password is required", op)
	}

	// Создание пользователя
	userID, err := s.userService.RegisterUser(ctx, req.GetEmail(), req.GetUsername(), req.GetPassword())
	if err != nil {
		s.log.Error("failed to register user", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to register user: %w", op, err)
	}

	s.log.Info("user registered", slog.String("op", op), slog.Int64("user_id", userID))

	return &auth.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	const op = "grpc.Login"

	if req.GetEmail() == "" {
		return nil, fmt.Errorf("%s: email is required", op)
	}

	if req.GetPassword() == "" {
		return nil, fmt.Errorf("%s: password is required", op)
	}

	id, err := s.userService.LoginUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		s.log.Error("failed to login user", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to login user: %w", op, err)
	}

	return &auth.LoginResponse{
		AccessToken: "",
		UserId:      id,
	}, nil
}
