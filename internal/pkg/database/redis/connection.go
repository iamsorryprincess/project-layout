package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Connection struct {
	*redis.Client
}

func New(config Config) (*Connection, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Username: config.User,
		Password: config.Password,
		DB:       config.DB,
	})

	if err := client.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("redis connection error: %v", err)
	}

	return &Connection{
		Client: client,
	}, nil
}
