package main

import (
	"context"
	"habr/db/auth"
	"habr/internal/auth/app"
	"habr/internal/auth/app/repositories"
	"habr/internal/auth/app/services"
	"habr/internal/auth/config"
	"habr/internal/auth/logger"
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
	userService := services.NewUserService(userRepo)
	id, err := userService.RegisterUser(ctx, "example_email", "bogdan", "hashed_password_example")
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("id:", id)

	application := app.New()
	go func() {
		application.Start(cfg, log)

	}()
}
