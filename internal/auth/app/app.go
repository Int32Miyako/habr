package app

import (
	"habr/internal/auth/app/grpc"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
)

type App struct {
	GRPCServer *grpc.App
}

func New() *App {
	return &App{}
}

func (app *App) Start(cfg *config.Config, log *slog.Logger, userService *services.UserService) {
	grpcApp := grpc.New(log, cfg, userService)
	app.GRPCServer = grpcApp

	grpcApp.MustRun()
}

func (app *App) Stop() {
	if app.GRPCServer != nil {
		app.GRPCServer.Stop()
	}
}
