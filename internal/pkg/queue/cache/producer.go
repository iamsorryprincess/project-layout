package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue"
)

type Producer[TMessage any] struct {
	fileName string
	mu       sync.Mutex
	messages []TMessage

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

func (p *Producer[TMessage]) Produce(_ context.Context, messages ...TMessage) error {
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
			return err
		}

		p.logger.Debug().Msgf("%d messages successfully handled", len(fileData))
		removeFile(p.fileName, p.logger)
	}

	p.mu.Lock()
	now := time.Now()
	if len(p.messages) == 0 {
		p.mu.Unlock()
		return nil
	}

	if err = saveToFile(p.fileName, p.messages, p.logger); err != nil {
		p.mu.Unlock()
		return fmt.Errorf("failed save logs to file: %v", err)
	}

	messages := p.messages
	p.messages = nil
	fmt.Println(time.Since(now))
	p.mu.Unlock()

	if err = p.producer.Produce(ctx, messages...); err != nil {
		return err
	}

	removeFile(p.fileName, p.logger)
	return nil
}
