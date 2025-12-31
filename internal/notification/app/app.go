package app

import (
	"habr/internal/notification/app/grpc"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
)

type App struct {
	GRPCServer  *grpc.App
	cfg         *config.Config
	log         *slog.Logger
	userService services.EmailService
}

func New(cfg *config.Config, log *slog.Logger, userService services.EmailService) *App {
	return &App{
		cfg:         cfg,
		log:         log,
		userService: userService,
	}
}

func (app *App) Start() error {
	grpcApp := grpc.New(app.log, app.cfg, app.userService)
	app.GRPCServer = grpcApp

	err := grpcApp.Run()
	return err
}

func (app *App) Stop() {
	if app.GRPCServer != nil {
		app.GRPCServer.Stop()
	}
}
