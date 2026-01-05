package consumer

import "context"

type MessageConsumer interface {
	Subscribe(ctx context.Context, topic string, handler func(*Message) error) error
	Close() error
}

type Message struct {
	Key   string
	Value []byte
}
