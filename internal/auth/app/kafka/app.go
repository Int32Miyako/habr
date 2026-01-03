package kafka

import (
	"fmt"
	"habr/internal/auth/app/kafka/producer"
	"habr/internal/auth/config"
	"log/slog"
)

type App struct {
	RegistrationNotifier *producer.RegistrationNotifier
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	brokers := cfg.Kafka.Brokers
	topic := cfg.Kafka.Topic

	registrationNotifier, err := producer.NewRegistrationNotifier(brokers, topic, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer.NewRegistrationNotifier: %w", err)
	}

	return &App{
		RegistrationNotifier: registrationNotifier,
	}, nil
}

func (app App) Close() error {
	err := app.RegistrationNotifier.Close()
	if err != nil {
		return fmt.Errorf("failed to close RegistrationNotifier: %w", err)
	}

	return nil
}
