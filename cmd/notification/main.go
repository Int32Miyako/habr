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
	ctx := context.Background()

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

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	serverErrors := make(chan error, 2)

	go func() {
		serverErrors <- application.Start(appCtx)
	}()

	select {
	case sig := <-stop:
		log.Info("Received shutdown signal", "signal", sig.String())
	case err = <-serverErrors:
		log.Error("Server error, shutting down", "error", err)
	}
	appCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()

	if err = application.Stop(shutdownCtx); err != nil {
		log.Error("Failed to stop application gracefully", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("Gracefully stopped")
}
