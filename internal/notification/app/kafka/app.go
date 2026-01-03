package kafka

import (
	"fmt"
	"habr/internal/notification/app/kafka/consumer"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
)

type App struct {
	RegistrationConsumer *consumer.RegistrationNotifier
}

func New(cfg *config.Config, log *slog.Logger, emailService services.EmailService) (*App, error) {
	notifier, err := consumer.NewRegistrationNotifier(cfg, log, emailService)
	if err != nil {
		return nil, fmt.Errorf("kafka app start: %w", err)
	}

	app := &App{
		RegistrationConsumer: notifier,
	}

	return app, nil
}

func (a App) Stop() error {
	err := a.RegistrationConsumer.Close()
	if err != nil {
		return fmt.Errorf("kafka app stop: %w", err)
	}

	return nil
}
