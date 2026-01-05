package kafka

import (
	"context"
	"fmt"
	"habr/internal/notification/app/kafka/consumer"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
)

type App struct {
	RegistrationConsumer *consumer.RegistrationNotifier
	TopicConsumer        string
}

func New(cfg *config.Config, log *slog.Logger, emailService services.EmailService, topic string) (*App, error) {
	notifier, err := consumer.NewRegistrationNotifier(cfg, log, emailService)
	if err != nil {
		return nil, fmt.Errorf("kafka app start: %w", err)
	}

	app := &App{
		RegistrationConsumer: notifier,
		TopicConsumer:        topic,
	}

	return app, nil
}

// Run запускает Kafka consumer и блокируется до отмены контекста.
func (a App) Run(ctx context.Context) error {
	const op = "kafka.Run"

	err := a.RegistrationConsumer.Subscribe(ctx, a.TopicConsumer)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	<-ctx.Done()
	return ctx.Err()
}

func (a App) Stop(ctx context.Context) error {
	done := make(chan error, 1)

	go func() {
		done <- a.RegistrationConsumer.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("kafka app stop: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("kafka app stop: graceful shutdown timeout exceeded")
	}
}
