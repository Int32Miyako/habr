package client

import (
	"fmt"
	"habr/internal/notification/core/models"
	"log/slog"

	"github.com/IBM/sarama"
)

type RegistrationNotifier struct {
	producer sarama.SyncProducer
	topic    string
	log      *slog.Logger
}

func (r *RegistrationNotifier) SendMessage(
	topic string,
	message *models.Message,
) error {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.Key),
		Value: sarama.StringEncoder(message.Value),
	}

	partition, offset, err := r.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send failed: %w", err)
	}

	r.log.Info(
		"message sent",
		slog.String("topic", topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)

	return nil
}

func NewProducer(brokers []string, topic string, log *slog.Logger) (*RegistrationNotifier, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama.NewSyncProducer: %w", err)
	}

	return &RegistrationNotifier{
		producer: producer,
		topic:    topic,
		log:      log,
	}, nil
}

func (r *RegistrationNotifier) Close() error {
	err := r.producer.Close()
	if err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}

	return nil
}
