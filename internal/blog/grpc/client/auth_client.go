package client

import (
	"context"
	"habr/protos/gen/go/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	return c.client.Register(ctx, &auth.RegisterRequest{
		Email:    email,
		Username: username,
		Password: password,
	})
}

// Login выполняет вход пользователя и возвращает токен
func (c *AuthClient) Login(ctx context.Context, email, password string) (*auth.LoginResponse, error) {
	return c.client.Login(ctx, &auth.LoginRequest{
		Email:    email,
		Password: password,
	})
}

// GetClient возвращает нативный gRPC клиент для прямого использования
func (c *AuthClient) GetClient() auth.AuthClient {
	return c.client
}
