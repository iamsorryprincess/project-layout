package queue

import "context"

type MessageHandler[TMessage any] interface {
	Handle(ctx context.Context, messages []TMessage) error
}

type Producer[TMessage any] interface {
	Produce(ctx context.Context, messages ...TMessage) error
}

type Consumer[TMessage any] interface {
	Consume(ctx context.Context) ([]TMessage, int64, error)
}
