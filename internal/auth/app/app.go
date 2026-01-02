package app

import (
	"habr/internal/auth/app/grpc"
	"habr/internal/auth/app/kafka/producer"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"log/slog"
	"os"

	"github.com/IBM/sarama"
)

type App struct {
	GRPCServer *grpc.App
	Producer   sarama.SyncProducer
}

func New() *App {
	return &App{}
}

func (app *App) Start(cfg *config.Config, log *slog.Logger, userService *services.UserService) {
	grpcApp := grpc.New(log, cfg, userService)

	prod, err := producer.New(cfg.Kafka.Brokers, cfg.Kafka.Topic, log)
	if err != nil {
		log.Error("failed to create kafka producer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	app.GRPCServer = grpcApp
	app.Producer = prod

	grpcApp.MustRun()
}

func (app *App) Stop() {
	if app.GRPCServer != nil {
		app.GRPCServer.Stop()
	}
}
