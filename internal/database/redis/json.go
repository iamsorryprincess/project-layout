package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *Connection) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal data %v error: %w", data, err)
	}

	if err = c.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set %v error: %w", key, err)
	}

	return nil
}

func (c *Connection) GetJSON(ctx context.Context, key string, value interface{}) error {
	data, err := c.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return err
		}
		return fmt.Errorf("redis get %v error: %w", key, err)
	}

	if err = json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("json unmarshal %v error: %w", data, err)
	}

	return nil
}
