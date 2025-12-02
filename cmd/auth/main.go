package main

import (
	"context"
	"habr/db/auth"
	"habr/internal/auth/app"
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

	_ = database

	application := app.New()
	go func() {
		application.Start(cfg, log)

	}()
}
