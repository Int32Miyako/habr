package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"habr/internal/notification/app/kafka/client"
	"habr/internal/notification/config"
	consumerContract "habr/internal/notification/core/interfaces/kafka/client"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
)

type RegistrationNotifier struct {
	consumer     consumerContract.MessageConsumer
	log          *slog.Logger
	emailService services.EmailService
}

func NewRegistrationNotifier(cfg *config.Config, log *slog.Logger, emailService services.EmailService) (*RegistrationNotifier, error) {
	consumer, err := client.NewConsumerGroup(cfg.Kafka, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &RegistrationNotifier{
		consumer:     consumer,
		log:          log,
		emailService: emailService,
	}, nil
}

func (c *RegistrationNotifier) Subscribe(ctx context.Context, topic string) error {
	err := c.consumer.Subscribe(ctx, topic, c.handleMessage)
	if err != nil {
		c.log.Error("failed to subscribe to topic",
			slog.String("topic", topic),
			slog.String("error", err.Error()),
		)

		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	return nil
}

func (c *RegistrationNotifier) handleMessage(msg *consumerContract.Message) error {
	c.log.Info("received message from kafka",
		slog.String("key", msg.Key),
		slog.String("value", string(msg.Value)),
	)

	// Здесь можно добавить логику обработки сообщения
	// Например, десериализация и отправка email
	var emailData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &emailData); err != nil {
		c.log.Error("failed to unmarshal message",
			slog.String("error", err.Error()),
			slog.String("value", string(msg.Value)),
		)

		return err
	}

	c.log.Info("message processed successfully",
		slog.Any("data", emailData),
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
