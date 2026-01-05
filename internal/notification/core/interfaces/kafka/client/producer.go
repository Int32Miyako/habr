package client

import (
	"habr/internal/notification/core/models"
)

type MessageProducer interface {
	SendMessage(topic string, message *models.Message) error
	Close() error
}
