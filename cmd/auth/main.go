package main

import (
	"context"
	"fmt"
	"habr/db/auth"
	"habr/internal/auth/app"
	"habr/internal/auth/app/kafka/producer"
	"habr/internal/auth/app/repositories"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"habr/internal/auth/core/jwt"
	"habr/internal/auth/logger"
	"habr/internal/notification/core/models"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	database, err := auth.Initialize(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	log := logger.New()
	log.Info("Starting auth service")

	userRepo := repositories.NewUserRepository(database.Pool)
	jwtManager := jwt.NewJWTManager(cfg)
	userService := services.NewUserService(userRepo, jwtManager)

	application := app.New(cfg, log, userService)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	serverErrors := make(chan error, 2)

	application.Start(serverErrors)

	// Send test messages to Kafka
	go func() {
		var rn *producer.RegistrationNotifier
		rn, err = producer.NewRegistrationNotifier(cfg.Kafka.Brokers, cfg.Kafka.Topic, log)
		if err != nil {
			log.Error("failed to create registration notifier", "error", err)
			serverErrors <- err
			return
		}
		defer func() {
			if err = rn.Close(); err != nil {
				log.Error("failed to close registration notifier", "error", err)
			}
		}()

		for i := 0; i < 200; i++ {
			msg := models.Message{
				Key: fmt.Sprintf("%d", i),
				Value: []byte(fmt.Sprintf(
					`{
					"email": "newuser%d@example.com",
					"username": "newuser%d",
					"type": "registration"
					}`, i, i)),
			}

			err = rn.SendMessage(&msg)
			if err != nil {
				log.Error("failed to send message", "error", err)
				serverErrors <- err
				return
			}
		}
	}()

	select {
	case sig := <-stop:
		log.Info("Received shutdown signal", "signal", sig.String())
	case err = <-serverErrors:
		log.Error("Server error, shutting down", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	defer cancel()

	application.Stop(shutdownCtx)
	log.Info("Gracefully stopped")
}
