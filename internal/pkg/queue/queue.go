package queue

import "context"

type Producer[TMessage any] interface {
	Produce(ctx context.Context, messages ...TMessage) error
}
