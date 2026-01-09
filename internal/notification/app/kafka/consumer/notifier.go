package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"habr/internal/notification/core/events"
	consumerContract "habr/internal/notification/core/interfaces/kafka/client"
	"habr/internal/notification/core/interfaces/services"
	"habr/internal/notification/core/models"
	"log/slog"
)

type RegistrationNotifier struct {
	consumer     consumerContract.MessageConsumer
	log          *slog.Logger
	emailService services.EmailService
}

func NewRegistrationNotifier(consumer consumerContract.MessageConsumer, log *slog.Logger, emailService services.EmailService) *RegistrationNotifier {
	return &RegistrationNotifier{
		consumer:     consumer,
		log:          log,
		emailService: emailService,
	}
}

func (c *RegistrationNotifier) Subscribe(ctx context.Context, topics []string) error {
	err := c.consumer.Subscribe(ctx, topics, c.handleMessage)
	if err != nil {
		c.log.Error("failed to subscribe to topic",
			slog.Any("topic", topics),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("failed to subscribe to topic %s: %w", topics, err)
	}

	return nil
}

func (c *RegistrationNotifier) handleMessage(msg *models.Message) error {
	var event events.UserRegistered
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		c.log.Error("failed to unmarshal event", err.Error())
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	c.log.Info("received message from kafka",
		slog.String("key", msg.Key),
		slog.Any("event", event),
	)

	return nil
}

func (c *RegistrationNotifier) Close() error {
	c.log.Info("closing kafka consumer")

	err := c.consumer.Close()
	if err != nil {
		c.log.Error("failed to close kafka consumer",
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("failed to close kafka consumer: %w", err)
	}

	return nil
}
