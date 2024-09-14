package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/database/redis"
)

type Producer[TMessage any] struct {
	key  string
	conn *redis.Connection
}

func NewProducer[TMessage any](key string, conn *redis.Connection) *Producer[TMessage] {
	return &Producer[TMessage]{
		key:  key,
		conn: conn,
	}
}

func (p *Producer[TMessage]) Produce(ctx context.Context, message TMessage) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message for key:%s : %w", p.key, err)
	}

	if err = p.conn.RPush(ctx, p.key, data).Err(); err != nil {
		return fmt.Errorf("redis failed RPUSH message for key:%s : %w", p.key, err)
	}

	return nil
}
