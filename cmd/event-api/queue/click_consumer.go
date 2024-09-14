package queue

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/domain"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

type ClickConsumer struct {
	logger log.Logger
	conn   *clickhouse.Connection
}

func NewClickConsumer(logger log.Logger, conn *clickhouse.Connection) *ClickConsumer {
	return &ClickConsumer{
		logger: logger,
		conn:   conn,
	}
}

func (c *ClickConsumer) Consume(ctx context.Context, clicks []domain.Click) error {
	if err := c.conn.SendBatch(ctx, "insert into clicks", clicks); err != nil {
		return fmt.Errorf("failed to insert clicks: %w", err)
	}
	return nil
}
