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
	"sync"

	"github.com/IBM/sarama"
)

type RegistrationNotifier struct {
	consumer     consumerContract.MessageConsumer
	log          *slog.Logger
	emailService services.EmailService
	wg           sync.WaitGroup
}

func NewRegistrationNotifier(cfg *config.Config, log *slog.Logger, emailService services.EmailService) (*RegistrationNotifier, error) {
	c, err := client.NewConsumerGroup(cfg.Kafka, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &RegistrationNotifier{
		consumer:     c,
		log:          log,
		emailService: emailService,
	}, nil
}

func (c *RegistrationNotifier) consumePartition(ctx context.Context, pc sarama.PartitionConsumer, partition int32) {
	defer c.wg.Done()
	defer func() {
		if err := pc.Close(); err != nil {
			c.log.Error("failed to close partition consumer",
				slog.String("error", err.Error()),
				slog.Int("partition", int(partition)),
			)
		}
	}()

	for {
		select {
		case msg := <-pc.Messages():
			if msg != nil {
				c.handleMessage(msg, partition)
			}
		case err := <-pc.Errors():
			if err != nil {
				c.log.Error("partition consumer error",
					slog.String("error", err.Error()),
					slog.Int("partition", int(partition)),
				)
			}
		case <-ctx.Done():
			c.log.Info("stopping partition consumer",
				slog.Int("partition", int(partition)),
			)

			return
		}
	}
}

func (c *RegistrationNotifier) handleMessage(msg *sarama.ConsumerMessage, partition int32) {
	c.log.Info("received message from kafka",
		slog.String("topic", msg.Topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", msg.Offset),
		slog.String("key", string(msg.Key)),
	)

	// Здесь можно добавить логику обработки сообщения
	// Например, десериализация и отправка email
	var emailData map[string]interface{}
	if err := json.Unmarshal(msg.Value, &emailData); err != nil {
		c.log.Error("failed to unmarshal message",
			slog.String("error", err.Error()),
			slog.String("value", string(msg.Value)),
		)

		return
	}

	c.log.Info("message processed successfully",
		slog.Any("data", emailData),
	)
}

func (c *RegistrationNotifier) Close() error {
	c.log.Info("closing kafka consumer")
	c.wg.Wait()

	err := c.consumer.Close()
	if err != nil {
		c.log.Error("failed to close kafka consumer")
	}

	return nil
}
