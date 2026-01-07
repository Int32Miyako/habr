package client

import (
	"habr/internal/notification/core/models"
)

type MessageProducer interface {
	SendMessage(message *models.Message) error
	Close() error
}
