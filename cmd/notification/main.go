package main

import (
	db "habr/db/notification"
	"habr/internal/notification/app"
	"habr/internal/notification/app/grpc"
	"habr/internal/notification/app/kafka/consumer"
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

	database := db.MustInitialize(cfg)
	defer database.Close()

	log := logger.SetupLogger(cfg.Env)
	log.Info("Starting notification service")

	emailRepo := repositories.NewEmailRepository(database.Pool)
	emailService := services.NewEmailService(emailRepo)

	grpcApp := grpc.New(log, cfg, emailService)

	cons, err := consumer.New(cfg.Kafka.Brokers, cfg.Kafka.GroupID, cfg.Kafka.Topic, log, emailService)
	if err != nil {
		log.Error("failed to create kafka consumer", slog.Any("error", err))
		os.Exit(1)
	}

	application := app.New(grpcApp, cons)
	go func() {
		if err := application.Start(); err != nil {
			log.Error("app stopped with error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	waitForShutdown(application)
	log.Info("Gracefully stopped")
}

func waitForShutdown(application *app.App) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
}
