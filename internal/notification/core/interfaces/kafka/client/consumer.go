package client

import (
	"context"
	"habr/internal/notification/core/models"
)

type MessageConsumer interface {
	Subscribe(ctx context.Context, topic string, handler func(*models.Message) error) error
	Close() error
}
