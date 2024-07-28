package repository

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/domain"
)

type EventRepository struct {
	conn *clickhouse.Connection
}

func NewEventRepository(conn *clickhouse.Connection) *EventRepository {
	return &EventRepository{
		conn: conn,
	}
}

func (r *EventRepository) Handle(ctx context.Context, events []domain.Event) error {
	const query = "INSERT INTO events (created_at, ip, country_id, platform_id)"
	if err := r.conn.SendBatch(ctx, query, events); err != nil {
		return fmt.Errorf("failed to send events batch: %w", err)
	}
	return nil
}
