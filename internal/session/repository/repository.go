package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/domain"
)

// Repository - common repository for all services
type Repository struct {
	key  string
	ttl  time.Duration
	conn *redis.Connection
}

func NewRepository(key string, ttl time.Duration, conn *redis.Connection) *Repository {
	return &Repository{
		key:  key,
		ttl:  ttl,
		conn: conn,
	}
}

func (r *Repository) Set(ctx context.Context, session domain.Session) error {
	key := fmt.Sprintf("%s:%s", r.key, session.IP)
	if err := r.conn.SetJSON(ctx, key, session, r.ttl); err != nil {
		return fmt.Errorf("redis failed to set session: %w", err)
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ip string) (domain.Session, error) {
	key := fmt.Sprintf("%s:%s", r.key, ip)
	var result domain.Session
	if err := r.conn.GetJSON(ctx, key, &result); err != nil {
		if errors.Is(err, err) {
			return domain.Session{}, domain.ErrNotFound{Message: fmt.Sprintf("session with ip %s is not found", ip)}
		}
		return domain.Session{}, fmt.Errorf("redis failed to get ip %s session: %w", ip, err)
	}
	return result, nil
}

func (r *Repository) UpdateTTL(ctx context.Context, ip string) error {
	key := fmt.Sprintf("%s:%s", r.key, ip)
	if err := r.conn.Expire(ctx, key, r.ttl).Err(); err != nil {
		return fmt.Errorf("redis failed to update ttl session: %w", err)
	}
	return nil
}
