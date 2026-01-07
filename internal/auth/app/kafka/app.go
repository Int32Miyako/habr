package kafka

import (
	"fmt"
	"habr/internal/auth/app/kafka/producer"
	"habr/internal/auth/app/kafka/producer/client"
	"log/slog"
)

type App struct {
	RegistrationNotifier *producer.RegistrationNotifier
}

func New(prod client.MessageProducer, log *slog.Logger) *App {
	return &App{
		RegistrationNotifier: producer.NewRegistrationNotifier(prod, log),
	}
}

func (app App) Close() error {
	err := app.RegistrationNotifier.Close()
	if err != nil {
		return fmt.Errorf("failed to close RegistrationNotifier: %w", err)
	}

	return nil
}
