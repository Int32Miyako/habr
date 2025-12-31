package server

import (
	"habr/internal/notification/core/interfaces/services"
	"habr/protos/gen/go/notification"
	"log/slog"

	"google.golang.org/grpc"
)

type serverAPI struct {
	notification.UnimplementedNotificationServer
	emailService services.EmailService
	log          *slog.Logger
}

func Register(server *grpc.Server, emailService services.EmailService, log *slog.Logger) {
	notification.RegisterNotificationServer(server, &serverAPI{
		emailService: emailService,
		log:          log,
	})
}
