package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue"
)

type Consumer[TMessage any] struct {
	fileName string

	logger   log.Logger
	handler  queue.MessageHandler[TMessage]
	producer queue.Producer[TMessage]
	consumer queue.Consumer[TMessage]
}

func NewConsumer[TMessage any](
	key string,
	logger log.Logger,
	handler queue.MessageHandler[TMessage],
	producer queue.Producer[TMessage],
	consumer queue.Consumer[TMessage],
) *Consumer[TMessage] {
	return &Consumer[TMessage]{
		fileName: fmt.Sprintf("consumer.redis.logs.%s", key),
		logger:   newLogger(key, logger),
		handler:  handler,
		producer: producer,
		consumer: consumer,
	}
}

func (c *Consumer[TMessage]) Consume(ctx context.Context) error {
	count, err := c.consume(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	for count > 0 {
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

func (c *Consumer[TMessage]) consume(ctx context.Context) (int64, error) {
	fileData, err := c.readFromFile()
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

		c.removeFile()
	}

	result, count, err := c.consumer.Consume(ctx)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	if err = c.saveToFile(result); err != nil {
		c.sendDataBack(ctx, result)
		return 0, fmt.Errorf("failed save logs to file: %v", err)
	}

	if err = c.handler.Handle(ctx, result); err != nil {
		c.sendDataBack(ctx, result)
		return 0, err
	}

	c.removeFile()
	return count, nil
}

func (c *Consumer[TMessage]) sendDataBack(ctx context.Context, data []TMessage) {
	c.logger.Debug().Msgf("sending failed %d messages back", len(data))

	if err := c.producer.Produce(ctx, data...); err != nil {
		c.logger.Error().Msgf("failed send data back: %v", err)
		return
	}

	c.logger.Debug().Msgf("%d failed messages sent back", len(data))
	c.removeFile()
}

func (c *Consumer[TMessage]) readFromFile() ([]TMessage, error) {
	file, err := os.OpenFile(c.fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			c.logger.Error().Msgf("failed close logs file: %v", err)
		}
	}()

	var result []TMessage
	if err = json.NewDecoder(file).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Consumer[TMessage]) saveToFile(result []TMessage) error {
	file, err := os.OpenFile(c.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			c.logger.Error().Msgf("failed close logs file: %v", err)
		}
	}()

	if err = json.NewEncoder(file).Encode(result); err != nil {
		return err
	}

	c.logger.Debug().Msgf("saved %d consuming messages to logs file", len(result))
	return nil
}

func (c *Consumer[TMessage]) removeFile() {
	if err := os.Remove(c.fileName); err != nil {
		c.logger.Error().Msgf("failed to remove logs file: %v", err)
	}
}
