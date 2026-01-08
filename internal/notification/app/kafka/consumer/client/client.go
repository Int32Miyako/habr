package client

import (
	"context"
	"fmt"
	"habr/internal/notification/config"
	"habr/internal/notification/core/models"
	"log/slog"

	"github.com/IBM/sarama"
)

type KafkaConsumerClient struct {
	consumerGroup sarama.ConsumerGroup
	log           *slog.Logger
	kafkaConfig   *config.Kafka
}

type ConsumerHandler struct {
	handler func(*models.Message) error
	log     *slog.Logger
}

// Setup is called once when the consumer group is started.
func (ch *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called once when the consumer group is terminated.
func (ch *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim is called for each claim returned by the consumer group.
// Process each message from the claim and mark the message as processed.
func (ch *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		msg := &models.Message{
			Key:   string(message.Key),
			Value: message.Value,
		}
		if err := ch.handler(msg); err != nil {
			ch.log.Error("handler error", "err", err)
			continue
		}

		session.MarkMessage(message, "processed")
	}

	return nil
}

func (k *KafkaConsumerClient) Subscribe(ctx context.Context, topics []string, handler func(*models.Message) error) error {
	consumerHandler := &ConsumerHandler{
		handler: handler,
		log:     k.log,
	}

	k.Run(ctx)

	go func() {
		for {
			err := k.consumerGroup.Consume(ctx, topics, consumerHandler)
			if err != nil {
				k.log.Error("Error consuming messages", slog.String("error", err.Error()))
				return
			}

			if ctx.Err() != nil {
				k.log.Info("Consumer context cancelled, stopping consumption")
				return
			}
		}
	}()

	return nil
}

func (k *KafkaConsumerClient) Close() error {
	err := k.consumerGroup.Close()
	if err != nil {
		return fmt.Errorf("failed to close consumer group: %w", err)
	}

	return nil
}

func NewKafkaConsumerClient(cfg *config.Kafka, log *slog.Logger) (*KafkaConsumerClient, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V4_1_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaConsumerClient{
		consumerGroup: consumerGroup,
		log:           log,
		kafkaConfig:   cfg,
	}, nil
}

func (k *KafkaConsumerClient) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case err, ok := <-k.consumerGroup.Errors():
				if !ok {
					return
				}
				k.log.Error("kafka consumer error", slog.Any("error", err))

			case <-ctx.Done():
				return
			}
		}
	}()
}
