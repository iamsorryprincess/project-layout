package cache

import (
	"context"
	"sync"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue"
)

type Producer[TMessage any] struct {
	mu       sync.Mutex
	messages []TMessage

	logger   log.Logger
	producer queue.Producer[TMessage]
}

func NewProducer[TMessage any](logger log.Logger, producer queue.Producer[TMessage]) *Producer[TMessage] {
	return &Producer[TMessage]{
		logger:   logger,
		producer: producer,
	}
}

func (p *Producer[TMessage]) Produce(_ context.Context, messages ...TMessage) error {
	p.mu.Lock()
	p.messages = append(p.messages, messages...)
	p.mu.Unlock()
	return nil
}

func (p *Producer[TMessage]) Send(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.messages) == 0 {
		return nil
	}

	if err := p.producer.Produce(ctx, p.messages...); err != nil {
		return err
	}

	p.messages = p.messages[:0]
	return nil
}
