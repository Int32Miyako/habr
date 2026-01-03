package main

import (
	"context"
	db "habr/db/notification"
	"habr/internal/notification/app"
	"habr/internal/notification/app/grpc"
	"habr/internal/notification/app/kafka"
	"habr/internal/notification/app/repositories"
	"habr/internal/notification/app/services"
	"habr/internal/notification/config"
	"habr/internal/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	ctx, cancel := context.WithCancel(context.Background())

	database := db.Initialize(ctx, cfg)
	defer database.Close()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting notification service")

	emailRepo := repositories.NewEmailRepository(database.Pool)
	emailService := services.NewEmailService(emailRepo)

	grpcApp := grpc.New(log, cfg, emailService)

	kafkaApp, err := kafka.New(cfg, log, emailService)
	if err != nil {
		log.Error("failed to create kafka consumer", slog.Any("error", err))
		os.Exit(1)
	}

	application := app.New(grpcApp, kafkaApp)

	go func() {
		if err = application.Start(ctx); err != nil {
			log.Error("app stopped with error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	cancel()

	err = application.Stop()
	if err != nil {
		log.Error("app stopped with error", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("Gracefully stopped")
}
