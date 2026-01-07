package producer

import (
	"fmt"
	"habr/internal/auth/app/kafka/producer/client"

	"habr/internal/auth/core/models"
	"log/slog"
)

type RegistrationNotifier struct {
	producer client.MessageProducer
	log      *slog.Logger
}

func NewRegistrationNotifier(producer client.MessageProducer, log *slog.Logger) *RegistrationNotifier {
	return &RegistrationNotifier{
		producer: producer,
		log:      log,
	}
}

func (p *RegistrationNotifier) SendMessage(message *models.Message) error {
	err := p.producer.SendMessage(message)
	if err != nil {
		p.log.Error("failed to send message to kafka",
			slog.String("error", err.Error()),
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
