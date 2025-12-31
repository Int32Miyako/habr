package server

import (
	"habr/internal/notification/core/interfaces/services"
	"habr/protos/gen/go/notification"
	"log/slog"

	"google.golang.org/grpc"
)

type serverAPI struct {
	notification.UnimplementedNotificationServer
	userService services.EmailService
	log         *slog.Logger
}

func Register(server *grpc.Server, service services.EmailService, log *slog.Logger) {
	notification.RegisterNotificationServer(server, &serverAPI{
		userService: service,
		log:         log,
	})
}
