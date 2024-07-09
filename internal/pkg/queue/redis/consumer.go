package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	redisdb "github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/redis/go-redis/v9"
)

type Consumer[TMessage any] struct {
	key   string
	count int

	logger log.Logger
	conn   *redisdb.Connection
}

func NewConsumer[TMessage any](key string, count int, logger log.Logger, conn *redisdb.Connection) *Consumer[TMessage] {
	return &Consumer[TMessage]{
		key:    key,
		count:  count,
		logger: newLogger(key, logger),
		conn:   conn,
	}
}

func (c *Consumer[TMessage]) Consume(ctx context.Context) ([]TMessage, int64, error) {
	result, err := c.conn.LPopCount(ctx, c.key, c.count).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, 0, fmt.Errorf("redis LPOP error: %v", err)
	}

	if len(result) == 0 {
		return nil, 0, nil
	}

	c.logger.Debug().Msgf("redis LPOP %d messages from list", len(result))

	messages := make([]TMessage, len(result))
	for i, data := range result {
		if err = json.Unmarshal([]byte(data), &messages[i]); err != nil {
			return nil, 0, fmt.Errorf("failed unmarshalling message: %v", err)
		}
	}

	count, err := c.conn.LLen(ctx, c.key).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("redis LLEN error: %v", err)
	}

	return messages, count, nil
}
