package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	*redis.Client

	logger log.Logger
}

func New(config Config, logger log.Logger) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Username: config.User,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection error: %v", err)
	}

	return &Connection{
		Client: client,
		logger: logger,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Client.Close(); err != nil {
		c.logger.Error().Str("type", "redis").Msg("redis failed to close connection")
	}
}

func (c *Connection) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal data %v error: %v", data, err)
	}

	if err = c.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set %v error: %v", key, err)
	}

	return nil
}

func (c *Connection) GetJSON(ctx context.Context, key string, value interface{}) error {
	data, err := c.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return err
		}
		return fmt.Errorf("redis get %v error: %v", key, err)
	}

	if err = json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("json unmarshal %v error: %v", data, err)
	}

	return nil
}
