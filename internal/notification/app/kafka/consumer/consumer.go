package consumer

import (
	"habr/internal/notification/core/interfaces/services"
	"log/slog"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
	log      *slog.Logger
}

func New(brokers []string, groupID string, topic string, log *slog.Logger, emailService services.EmailService) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func (p *Consumer) Close() error {
	return p.consumer.Close()
}
