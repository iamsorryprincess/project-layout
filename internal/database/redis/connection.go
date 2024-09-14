package redis

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	logger log.Logger
	*redis.Client
}

func New(logger log.Logger, config Config) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Username: config.User,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection error: %w", err)
	}

	return &Connection{
		logger: logger,
		Client: client,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Client.Close(); err != nil {
		c.logger.Error().Err(err).Msg("redis failed to close connection")
	}
}
