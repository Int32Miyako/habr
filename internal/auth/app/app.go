package app

import (
	"habr/internal/auth/app/grpc"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"

	"google.golang.org/grpc/serviceconfig"
)

type App struct {
	// GRPCServer *grpcapp.App
}

func New() *App {
	return &App{}
}

func (app *App) Start(cfg *config.Config, log *slog.Logger, userService *services.UserService) {
	grpcApp := grpc.New(log, cfg, userService)

}
