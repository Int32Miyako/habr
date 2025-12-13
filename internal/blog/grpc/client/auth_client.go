package client

import (
	"habr/protos/gen/go/auth"

	"google.golang.org/grpc"
)

type AuthClient struct {
	client auth.AuthClient
	conn   *grpc.ClientConn
}

func NewAuthClient(authAddr string) (*AuthClient, error) {
	grpcConn, err := grpc.NewClient(authAddr)
	if err != nil {
		return nil, err
	}

	defer grpcConn.Close()
	client := auth.NewAuthClient(grpcConn)

	return &AuthClient{
		client: client,
		conn:   grpcConn,
	}, nil
}
