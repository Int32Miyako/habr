package kafka

import (
	"fmt"
	"habr/internal/auth/app/kafka/producer"
	"habr/internal/auth/core/models"
	"log/slog"
)

type App struct {
	RegistrationNotifier *producer.RegistrationNotifier
}

type messageProducer interface {
	SendMessage(message *models.Message) error
	Close() error
}

func New(messageProducer messageProducer, log *slog.Logger) *App {
	return &App{
		RegistrationNotifier: producer.NewRegistrationNotifier(messageProducer, log),
	}
}

func (app App) Close() error {
	err := app.RegistrationNotifier.Close()
	if err != nil {
		return fmt.Errorf("failed to close RegistrationNotifier: %w", err)
	}

	return nil
}
