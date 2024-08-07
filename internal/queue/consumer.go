package queue

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/iamsorryprincess/project-layout/internal/log"
)

type MessageHandler[TMessage any] interface {
	Handle(ctx context.Context, messages []TMessage) error
}

type Consumer[TMessage any] interface {
	Consume(ctx context.Context) ([]TMessage, int64, error)
}

type FileCachingConsumer[TMessage any] struct {
	fileName string
	count    int64

	logger   log.Logger
	handler  MessageHandler[TMessage]
	producer Producer[TMessage]
	consumer Consumer[TMessage]
}

func NewFileCachingConsumer[TMessage any](
	key string,
	count int,
	logger log.Logger,
	handler MessageHandler[TMessage],
	producer Producer[TMessage],
	consumer Consumer[TMessage],
) *FileCachingConsumer[TMessage] {
	return &FileCachingConsumer[TMessage]{
		fileName: fmt.Sprintf("consumer.logs.%s", key),
		count:    int64(count),
		logger:   newLogger(key, "consumer", logger),
		handler:  handler,
		producer: producer,
		consumer: consumer,
	}
}

func (c *FileCachingConsumer[TMessage]) Consume(ctx context.Context) error {
	count, err := c.consume(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	for count >= c.count {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			count, err = c.consume(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *FileCachingConsumer[TMessage]) consume(ctx context.Context) (int64, error) {
	fileData, err := readFromFile[TMessage](c.fileName, c.logger)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("failed open logs file: %v", err)
		}
	} else {
		c.logger.Debug().Msgf("read %d messages from local file", len(fileData))

		if err = c.handler.Handle(ctx, fileData); err != nil {
			return 0, fmt.Errorf("failed handle messages: %v", err)
		}

		c.logger.Debug().Msgf("%d messages successfully handled", len(fileData))
		removeFile(c.fileName, c.logger)
	}

	result, count, err := c.consumer.Consume(ctx)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	if err = saveToFile(c.fileName, result, c.logger); err != nil {
		c.sendDataBack(ctx, result)
		return 0, fmt.Errorf("failed save logs to file: %v", err)
	}

	c.logger.Debug().Msgf("saved %d consuming messages to logs file", len(result))

	if err = c.handler.Handle(ctx, result); err != nil {
		c.sendDataBack(ctx, result)
		return 0, err
	}

	c.logger.Debug().Msgf("%d messages successfully handled", len(result))
	removeFile(c.fileName, c.logger)
	return count, nil
}

func (c *FileCachingConsumer[TMessage]) sendDataBack(ctx context.Context, data []TMessage) {
	c.logger.Debug().Msgf("sending failed %d messages back", len(data))

	if err := c.producer.Produce(ctx, data...); err != nil {
		c.logger.Error().Msgf("failed send data back: %v", err)
		return
	}

	c.logger.Debug().Msgf("%d failed messages sent back", len(data))
	removeFile(c.fileName, c.logger)
}
