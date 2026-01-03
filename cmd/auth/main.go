package main

import (
	"context"
	"habr/db/auth"
	"habr/internal/auth/app"
	"habr/internal/auth/app/repositories"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"habr/internal/auth/core/jwt"
	"habr/internal/auth/logger"
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

	log := logger.New()
	log.Info("Starting auth service")

	userRepo := repositories.NewUserRepository(database.Pool)
	jwtManager := jwt.NewJWTManager(cfg)
	userService := services.NewUserService(userRepo, jwtManager)

	application := app.New(cfg, log, userService)

	go func() {
		application.Start()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.Stop()
	log.Info("Gracefully stopped")
}
