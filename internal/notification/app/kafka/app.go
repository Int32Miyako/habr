package kafka

import (
	"context"
	"fmt"
	"habr/internal/notification/app/kafka/consumer"
	"habr/internal/notification/app/kafka/consumer/client"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
)

type App struct {
	RegistrationConsumer *consumer.RegistrationNotifier
	TopicsConsumer       []string
}

func New(cfg *config.Config, log *slog.Logger, emailService services.EmailService) (*App, error) {
	cons, err := client.NewKafkaConsumerClient(cfg.Kafka, log)
	if err != nil {
		return nil, fmt.Errorf("kafka app start: %w", err)
	}

	notifier := consumer.NewRegistrationNotifier(cons, log, emailService)

	app := &App{
		RegistrationConsumer: notifier,
		TopicsConsumer:       cfg.Kafka.Topics,
	}

	return app, nil
}

// Run запускает Kafka consumer и блокируется до отмены контекста.
func (a App) Run(ctx context.Context) error {
	const op = "kafka.Run"

	err := a.RegistrationConsumer.Subscribe(ctx, a.TopicsConsumer)
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
