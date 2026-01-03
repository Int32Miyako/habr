package producer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

type RegistrationNotifier struct {
	producer sarama.SyncProducer
	topic    string
	log      *slog.Logger
}

func NewRegistrationNotifier(brokers []string, topic string, log *slog.Logger) (*RegistrationNotifier, error) {
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

func (p *RegistrationNotifier) SendMessage(key string, message interface{}) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		p.log.Error("failed to marshal message", slog.String("error", err.Error()))
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.log.Error("failed to send message to kafka",
			slog.String("error", err.Error()),
			slog.String("topic", p.topic),
		)

		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	p.log.Info("message sent to kafka",
		slog.String("topic", p.topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)

	return nil
}

func (p *RegistrationNotifier) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("close kafka producer: %w", err)
	}

	return nil
}
