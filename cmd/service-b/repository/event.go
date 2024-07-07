package repository

import (
	"context"

	"github.com/iamsorryprincess/project-layout/internal/app/domain"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
)

type EventRepository struct {
	conn *clickhouse.Connection
}

func NewEventRepository(conn *clickhouse.Connection) *EventRepository {
	return &EventRepository{
		conn: conn,
	}
}

func (r *EventRepository) Save(ctx context.Context, events []domain.Event) error {
	return nil
}
