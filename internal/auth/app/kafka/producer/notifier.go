package producer

import (
	"fmt"
	"habr/internal/auth/app/kafka/producer/client"
	producerContract "habr/internal/notification/core/interfaces/kafka/client"
	"habr/internal/notification/core/models"
	"log/slog"
)

type RegistrationNotifier struct {
	producer producerContract.MessageProducer
	topic    string
	log      *slog.Logger
}

func NewRegistrationNotifier(brokers []string, topic string, log *slog.Logger) (*RegistrationNotifier, error) {
	producer, err := client.NewProducer(brokers, topic, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &RegistrationNotifier{
		producer: producer,
		topic:    topic,
		log:      log,
	}, nil
}

func (p *RegistrationNotifier) SendMessage(message *models.Message) error {
	err := p.producer.SendMessage(p.topic, message)
	if err != nil {
		p.log.Error("failed to send message to kafka",
			slog.String("error", err.Error()),
			slog.String("topic", p.topic),
		)

		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}

func (p *RegistrationNotifier) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("close kafka producer: %w", err)
	}

	return nil
}
