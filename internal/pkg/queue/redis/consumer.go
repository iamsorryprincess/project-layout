package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	redisdb "github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/redis/go-redis/v9"
)

type ConsumeHandler[TMessage any] func(ctx context.Context, messages []TMessage) error

type Consumer[TMessage any] struct {
	fileName string

	key   string
	count int

	logger log.Logger

	conn *redisdb.Connection

	handlerFunc ConsumeHandler[TMessage]
}

func NewConsumer[TMessage any](key string, count int, logger log.Logger, conn *redisdb.Connection, handlerFunc ConsumeHandler[TMessage]) *Consumer[TMessage] {
	return &Consumer[TMessage]{
		fileName:    fmt.Sprintf("consumer.redis.logs.%s", key),
		key:         key,
		count:       count,
		logger:      newLogger(key, logger),
		conn:        conn,
		handlerFunc: handlerFunc,
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
		c.logger.Debug().Msg("reading messages from local logs file")
		messages := make([]TMessage, len(fileData))
		for i, data := range fileData {
			if err = json.Unmarshal([]byte(data), &messages[i]); err != nil {
				return 0, fmt.Errorf("failed unmarshalling message: %v", err)
			}
		}

		c.logger.Debug().Msgf("read %d messages from local file", len(messages))

		if err = c.handlerFunc(ctx, messages); err != nil {
			return 0, fmt.Errorf("failed handle messages: %v", err)
		}

		c.logger.Debug().Msgf("%d messages successfully handled", len(messages))

		c.removeFile()
	}

	result, err := c.conn.LPopCount(ctx, c.key, c.count).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("redis LPOP error: %v", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	if err = c.saveToFile(result); err != nil {
		c.sendDataBack(ctx, result)
		return 0, fmt.Errorf("failed save logs to file: %v", err)
	}

	messages := make([]TMessage, len(result))
	for i, data := range result {
		if err = json.Unmarshal([]byte(data), &messages[i]); err != nil {
			c.sendDataBack(ctx, result)
			return 0, fmt.Errorf("failed unmarshalling message: %v", err)
		}
	}

	if err = c.handlerFunc(ctx, messages); err != nil {
		c.sendDataBack(ctx, result)
		return 0, fmt.Errorf("failed handle messages: %v", err)
	}

	c.logger.Debug().Msgf("%d messages successfully handled", len(messages))
	c.removeFile()

	count, err := c.conn.LLen(ctx, c.key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis LLEN error: %v", err)
	}

	return count, nil
}

func (c *Consumer[TMessage]) sendDataBack(ctx context.Context, data []string) {
	interfaceData := make([]interface{}, len(data))
	for i := range data {
		interfaceData[i] = data[i]
	}

	c.logger.Debug().Msgf("sending failed %d messages back", len(data))

	if err := c.conn.RPush(ctx, c.key, interfaceData...).Err(); err != nil {
		c.logger.Error().Msgf("failed RPUSH data back: %v", err)
		return
	}

	c.logger.Debug().Msgf("sending failed %d messages back", len(data))
	c.removeFile()
}

func (c *Consumer[TMessage]) saveToFile(result []string) error {
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

func (c *Consumer[TMessage]) readFromFile() ([]string, error) {
	file, err := os.OpenFile(c.fileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := file.Close(); cErr != nil {
			c.logger.Error().Msgf("failed close logs file: %v", err)
		}
	}()

	var result []string
	if err = json.NewDecoder(file).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Consumer[TMessage]) removeFile() {
	if err := os.Remove(c.fileName); err != nil {
		c.logger.Error().Msgf("failed to remove logs file: %v", err)
	}
}
