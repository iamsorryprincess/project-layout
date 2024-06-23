package redis

import (
	"context"
	"fmt"

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
