package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue"
)

type Producer[TMessage any] struct {
	mu       sync.Mutex
	messages []TMessage

	fileName string
	logger   log.Logger
	producer queue.Producer[TMessage]
}

func NewProducer[TMessage any](key string, logger log.Logger, producer queue.Producer[TMessage]) *Producer[TMessage] {
	return &Producer[TMessage]{
		fileName: fmt.Sprintf("producer.logs.%s", key),
		logger:   newLogger(key, "producer", logger),
		producer: producer,
	}
}

func (p *Producer[TMessage]) Produce(ctx context.Context, messages ...TMessage) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	p.messages = append(p.messages, messages...)
	p.mu.Unlock()
	return nil
}

func (p *Producer[TMessage]) Send(ctx context.Context) error {
	fileData, err := readFromFile[TMessage](p.fileName, p.logger)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed open logs file: %v", err)
		}
	} else {
		p.logger.Debug().Msgf("read %d messages from local file", len(fileData))

		if err = p.producer.Produce(ctx, fileData...); err != nil {
			return fmt.Errorf("failed to produce messages: %v", err)
		}

		p.logger.Debug().Msgf("%d messages from logs file successfully sent", len(fileData))
		removeFile(p.fileName, p.logger)
	}

	p.mu.Lock()
	messages := p.messages
	p.messages = nil
	p.mu.Unlock()

	if len(messages) == 0 {
		return nil
	}

	if err = p.producer.Produce(ctx, messages...); err != nil {
		p.logger.Error().Msgf("failed to send messages; saving messages to logs file: %v", err)
		if fErr := saveToFile[TMessage](p.fileName, messages, p.logger); fErr != nil {
			p.mu.Lock()
			p.messages = append(p.messages, messages...)
			p.mu.Unlock()
			return fmt.Errorf("failed to produce messages: %v; failed to save messages to logs file: %v", err, fErr)
		}

		p.logger.Debug().Msgf("%d failed messages saved to logs file", len(messages))
		return fmt.Errorf("failed to produce messages: %v", err)
	}

	return nil
}
