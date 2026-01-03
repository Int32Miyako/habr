package app

import (
	"habr/internal/auth/app/grpc"
	"habr/internal/auth/app/kafka"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
	"os"
)

type App struct {
	GRPCApp  *grpc.App
	KafkaApp *kafka.App
	log      *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger, userService *services.UserService) *App {
	grpcApp := grpc.New(log, cfg, userService)

	kafkaApp, err := kafka.New(cfg, log)
	if err != nil {
		log.Error("failed to create kafka producer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return &App{
		GRPCApp:  grpcApp,
		KafkaApp: kafkaApp,
		log:      log,
	}
}

func (app *App) Start() {
	app.GRPCApp.MustRun()
}

func (app *App) Stop() {
	if app.GRPCApp != nil {
		app.GRPCApp.Stop()
	}

	if app.KafkaApp != nil {
		err := app.KafkaApp.Close()
		if err != nil {
			app.log.Error("failed to close kafka prod", slog.String("error", err.Error()))
		}
	}
}
