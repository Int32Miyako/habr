package main

import (
	db "habr/db/notification"
	"habr/internal/notification/app"
	"habr/internal/notification/config"
	"habr/internal/notification/repositories"
	"habr/internal/notification/services"
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

	application := app.New(cfg, log, emailService)
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
