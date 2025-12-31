package main

import (
	db "habr/db/notification"
	"habr/internal/notification/config"
	"habr/internal/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	database := db.MustInitialize(cfg)
	_ = database
	log := logger.New()
	log.Info("Starting notification service")
}
