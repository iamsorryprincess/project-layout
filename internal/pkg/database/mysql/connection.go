package mysql

import (
	"database/sql"
	"fmt"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

type Connection struct {
	*sql.DB

	logger log.Logger
}

func New(config Config, logger log.Logger) (*Connection, error) {
	db, err := sql.Open("mysql", config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("mysql failed to connect: %w", err)
	}

	db.SetConnMaxLifetime(config.ConnectionMaxLifetime.Duration)
	if config.ConnectionMaxIdleTime.Duration > time.Second*0 {
		db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime.Duration)
	}
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("mysql failed to ping: %w", err)
	}

	return &Connection{
		DB:     db,
		logger: logger,
	}, nil
}

func (c *Connection) Close() {
	if err := c.DB.Close(); err != nil {
		c.logger.Error().Str("type", "mysql").Msg("mysql failed to close connection")
	}
}
