package redis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	redisdatabase "github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/iamsorryprincess/project-layout/internal/queue"
	"github.com/redis/go-redis/v9"
)

type ConsumerConfig struct {
	Key             string
	Count           int
	WorkersCount    int
	ConsumeInterval time.Duration
}

type ConsumerWorker[TMessage any] struct {
	logger log.Logger

	name   string
	config ConsumerConfig

	conn     *redisdatabase.Connection
	consumer queue.Consumer[TMessage]

	wg sync.WaitGroup
}

func NewConsumerWorker[TMessage any](logger log.Logger, name string, config ConsumerConfig, conn *redisdatabase.Connection, consumer queue.Consumer[TMessage]) *ConsumerWorker[TMessage] {
	return &ConsumerWorker[TMessage]{
		logger:   logger,
		name:     name,
		config:   config,
		conn:     conn,
		consumer: consumer,
		wg:       sync.WaitGroup{},
	}
}

func (c *ConsumerWorker[TMessage]) Start(ctx context.Context) {
	for i := 0; i < c.config.WorkersCount; i++ {
		c.wg.Add(1)
		go func(workerID int) {
			defer c.wg.Done()

			timer := time.NewTimer(c.config.ConsumeInterval)
			defer timer.Stop()

			for {
				select {
				case <-ctx.Done():
					c.logger.Info().
						Str("worker_name", c.name).
						Int("worker_id", workerID).
						Msg("consumer worker stopped due to context cancellation")
					return
				case <-timer.C:
					if err := c.processMessages(ctx); err != nil {
						if !errors.Is(err, context.Canceled) {
							c.logger.Error().
								Str("worker_name", c.name).
								Int("worker_id", workerID).
								Err(err).
								Msg("error while consumer processing messages")
						}
					}
					timer.Reset(c.config.ConsumeInterval)
				}
			}
		}(i)
	}
}

func (c *ConsumerWorker[TMessage]) processMessages(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		data, err := c.conn.LPopCount(ctx, c.config.Key, c.config.Count).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil
			}

			return err
		}

		if len(data) == 0 {
			return nil
		}

		var buf bytes.Buffer
		buf.WriteString("[")

		n := len(data)
		for i, value := range data {
			buf.WriteString(value)

			if i != n-1 {
				buf.WriteString(",")
			}
		}

		buf.WriteString("]")

		var messages []TMessage
		if err = json.Unmarshal(buf.Bytes(), &messages); err != nil {
			return err
		}

		if err = c.consumer.Consume(ctx, messages); err != nil {
			return err
		}
	}
}

func (c *ConsumerWorker[TMessage]) Shutdown() {
	c.wg.Wait()
	c.logger.Info().Msg("all consuming workers stopped")
}
