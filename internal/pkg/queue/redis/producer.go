package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
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

func (p *Producer[TMessage]) Produce(ctx context.Context, messages ...TMessage) error {
	switch len(messages) {
	case 0:
		return nil
	case 1:
		data, err := json.Marshal(messages[0])
		if err != nil {
			return fmt.Errorf("failed to marshal message %v with key: %s: %w", messages[0], p.key, err)
		}

		if err = p.conn.RPush(ctx, p.key, data).Err(); err != nil {
			return fmt.Errorf("failed RPush message %v with key: %s: %w", messages[0], p.key, err)
		}

		return nil
	default:
		messagesData := make([]interface{}, len(messages))
		for i, message := range messages {
			data, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("failed to marshal message %v with key: %s: %w", messages[i], p.key, err)
			}
			messagesData[i] = data
		}

		if err := p.conn.RPush(ctx, p.key, messagesData...).Err(); err != nil {
			return fmt.Errorf("failed RPush messages with key: %s: %w", p.key, err)
		}

		return nil
	}
}
