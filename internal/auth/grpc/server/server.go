package server

import (
	"context"
	"fmt"
	"habr/internal/auth/app/services"
	"habr/internal/auth/core/constants"
	"habr/internal/blog/http-server/dto"
	"habr/protos/gen/go/auth"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.Internal, constants.ErrInternalServer.Error())
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

	user, err := s.userService.LoginUser(ctx, dto.RequestLoginUser{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})

	if err != nil {
		s.log.Error("failed to login user", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to login user: %w", op, err)
	}

	return &auth.LoginResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		UserId:       user.UserId,
	}, nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *auth.RefreshRequest) (*auth.RefreshResponse, error) {
	const op = "grpc.Refresh"

	if req.GetRefreshToken() == "" {
		return nil, fmt.Errorf("%s: refresh token is required", op)
	}

	accessToken, err := s.userService.RefreshTokens(ctx, req.GetRefreshToken())
	if err != nil {
		s.log.Error("failed to refresh tokens", slog.String("op", op), slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: failed to refresh tokens: %w", op, err)
	}

	s.log.Info("tokens refreshed", slog.String("op", op))

	return &auth.RefreshResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *serverAPI) Validate(ctx context.Context, req *auth.ValidateRequest) (*auth.ValidateResponse, error) {
	const op = "grpc.Validate"

	if req.GetAccessToken() == "" {
		return nil, fmt.Errorf("%s: access token is required", op)
	}

	claims, err := s.userService.ValidateAccessToken(ctx, req.GetAccessToken())
	if err != nil {
		s.log.Error("failed to validate access token", slog.String("op", op), slog.String("error", err.Error()))
		return &auth.ValidateResponse{
			Valid:  false,
			UserId: 0,
		}, nil
	}

	s.log.Info("access token validated", slog.String("op", op), slog.Int64("user_id", claims.UserID))

	return &auth.ValidateResponse{
		Valid:  true,
		UserId: claims.UserID,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	const op = "grpc.Logout"

	if req.GetAccessToken() == "" {
		return nil, fmt.Errorf("%s: access token is required", op)
	}

	s.log.Info("user logged out", slog.String("op", op))

	return &auth.LogoutResponse{
		Success: true,
	}, nil
}
