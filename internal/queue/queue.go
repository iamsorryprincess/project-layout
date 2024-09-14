package queue

import "context"

type Producer[TMessage any] interface {
	Produce(ctx context.Context, message TMessage) error
}

type Consumer[TMessage any] interface {
	Consume(ctx context.Context, messages []TMessage) error
}
