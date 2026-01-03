package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/services"
	"log/slog"
	"sync"

	"github.com/IBM/sarama"
)

type RegistrationNotifier struct {
	consumer     sarama.Consumer
	topic        string
	log          *slog.Logger
	emailService services.EmailService
	wg           sync.WaitGroup
}

func NewRegistrationNotifier(cfg *config.Config, log *slog.Logger, emailService services.EmailService) (*RegistrationNotifier, error) {
	brokers := cfg.Kafka.Brokers
	topic := cfg.Kafka.Topic

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V4_1_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &RegistrationNotifier{
		consumer:     consumer,
		topic:        topic,
		log:          log,
		emailService: emailService,
	}, nil
}

func (c *RegistrationNotifier) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	c.log.Info("starting kafka consumer",
		slog.String("topic", c.topic),
		slog.Int("partitions", len(partitions)),
	)

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("failed to consume partition %d: %w", partition, err)
		}

		c.wg.Add(1)

		go c.consumePartition(ctx, pc, partition)
	}

	return nil
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
