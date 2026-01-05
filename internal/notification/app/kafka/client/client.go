package client

import (
	"fmt"
	"habr/internal/notification/config"

	"github.com/IBM/sarama"
)

func NewSaramaConfig() *sarama.Config {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V4_1_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	return saramaCfg
}

func NewConsumer(cfg *config.Config) (sarama.Consumer, error) {
	saramaCfg := NewSaramaConfig()
	consumer, err := sarama.NewConsumer(cfg.Kafka.Brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}
	return consumer, nil
}
