package client

import (
	"context"
	"fmt"
	"habr/internal/notification/config"
	"habr/internal/notification/core/interfaces/kafka/client"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
)

type ConsumerHandler struct {
	handler func(*client.Message) error
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
		msg := &client.Message{
			Key:   string(message.Key),
			Value: message.Value,
		}
		if err := ch.handler(msg); err != nil {
			ch.log.Error("handler error", "err", err)
		}
		session.MarkMessage(message, "processed")
	}

	return nil
}

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	log           *slog.Logger
	kafkaConfig   *config.Kafka
}

func (k *KafkaConsumer) Subscribe(ctx context.Context, topic string, handler func(*client.Message) error) error {
	consumerHandler := &ConsumerHandler{
		handler: handler,
		log:     k.log,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				k.log.Info("Consumer context cancelled, stopping consumption")
				return
			default:
			}

			err := k.consumerGroup.Consume(ctx, strings.Split(topic, ","), consumerHandler)
			if err != nil {
				k.log.Error("Error consuming messages", slog.String("error", err.Error()))
				return
			}
		}
	}()

	return nil
}

func (k *KafkaConsumer) Close() error {
	err := k.consumerGroup.Close()
	if err != nil {
		return err
	}

	return nil
}

func NewConsumerGroup(cfg *config.Kafka, log *slog.Logger) (*KafkaConsumer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V4_1_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.ConsumerGroup, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		log:           log,
		kafkaConfig:   cfg,
	}, nil
}
