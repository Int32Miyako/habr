package main

import (
	"habr/db/notification"
	"habr/internal/notification/config"
	"habr/internal/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	notification.MustInitialize(cfg)

	log := logger.New()
	log.Info("Starting notification service")
}
