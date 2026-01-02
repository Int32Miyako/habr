package producer

import (
	"log/slog"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	log      *slog.Logger
}

func New(brokers []string, topic string, log *slog.Logger) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
