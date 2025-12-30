package client

import (
	"context"
	"habr/internal/auth/core/constants"
	"habr/protos/gen/go/auth"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	client auth.AuthClient
	conn   *grpc.ClientConn
}

func NewAuthClient(authAddr string) (*AuthClient, error) {
	grpcConn, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthClient(grpcConn)

	return &AuthClient{
		client: client,
		conn:   grpcConn,
	}, nil
}

// Close закрывает соединение с gRPC сервером
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Register регистрирует нового пользователя
func (c *AuthClient) Register(ctx context.Context, email, username, password string) (*auth.RegisterResponse, error) {
	resp, err := c.client.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return nil, constants.ErrUserAlreadyExists
			case codes.InvalidArgument:
				return nil, constants.ErrInvalidCredentials
			default:
				return nil, constants.ErrInternalServer
			}
		}
		return nil, constants.ErrInternalServer
	}

	return resp, nil

}

// Login выполняет вход пользователя и возвращает токен
func (c *AuthClient) Login(ctx context.Context, email, password string) (*auth.LoginResponse, error) {
	resp, err := c.client.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, constants.ErrUserNotFound
			case codes.Unauthenticated:
				return nil, constants.ErrInvalidCredentials
			default:
				return nil, constants.ErrInternalServer
			}
		}
		return nil, constants.ErrInternalServer
	}

	return resp, nil

}

// GetClient возвращает нативный gRPC клиент для прямого использования
func (c *AuthClient) GetClient() auth.AuthClient {
	return c.client
}

func (c *AuthClient) Validate(ctx context.Context, token string) (bool, error) {
	resp, err := c.client.Validate(ctx, &auth.ValidateRequest{
		AccessToken: token,
	})

	if err != nil {
		log.Println("validate error:", err)
		return false, err
	}

	return resp.Valid, nil
}

func (c *AuthClient) Refresh(ctx context.Context, refreshToken string) (*auth.RefreshResponse, error) {
	return c.client.Refresh(ctx, &auth.RefreshRequest{
		RefreshToken: refreshToken,
	})
}
