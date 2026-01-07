package client

import (
	"fmt"
	"habr/internal/auth/core/models"
	"log/slog"

	"github.com/IBM/sarama"
)

type MessageProducer interface {
	SendMessage(message *models.Message) error
	Close() error
}

type ProducerKafkaClient struct {
	producer sarama.SyncProducer
	topic    string
	log      *slog.Logger
}

func (r *ProducerKafkaClient) SendMessage(
	message *models.Message,
) error {

	msg := &sarama.ProducerMessage{
		Topic: r.topic,
		Key:   sarama.StringEncoder(message.Key),
		Value: sarama.ByteEncoder(message.Value),
	}

	partition, offset, err := r.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("kafka send failed: %w", err)
	}

	r.log.Info(
		"message sent",
		slog.String("topic", r.topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)

	return nil
}

func NewProducerKafkaClient(brokers []string, topic string, log *slog.Logger) (*ProducerKafkaClient, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama.NewSyncProducer: %w", err)
	}

	return &ProducerKafkaClient{
		producer: producer,
		topic:    topic,
		log:      log,
	}, nil
}

func (r *ProducerKafkaClient) Close() error {
	err := r.producer.Close()
	if err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}

	return nil
}
