package grpc

import (
	"habr/protos/gen/go/auth"

	"google.golang.org/grpc"
)

type grpcService struct {
}

func NewGrpcService() *grpcService {
	return &grpcService{}
}

type serverAPI struct {
	auth.UnimplementedAuthServer
}

func Register(gRPCServer *grpc.Server) {

}
