package mysql

import (
	"database/sql"
	"fmt"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

type Connection struct {
	logger log.Logger
	*sql.DB
}

func New(logger log.Logger, config Config) (*Connection, error) {
	db, err := sql.Open("mysql", config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("mysql failed to connect: %w", err)
	}

	db.SetConnMaxLifetime(config.ConnectionMaxLifetime)
	if config.ConnectionMaxIdleTime > time.Second*0 {
		db.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	}
	db.SetMaxIdleConns(config.MaxIdleConnections)
	db.SetMaxOpenConns(config.MaxOpenConnections)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("mysql failed to ping: %w", err)
	}

	return &Connection{
		logger: logger,
		DB:     db,
	}, nil
}

func (c *Connection) Close() {
	if err := c.DB.Close(); err != nil {
		c.logger.Error().Err(err).Msg("mysql failed to close")
	}
}

func (c *Connection) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		c.logger.Error().Err(err).Msg("mysql failed to close rows")
	}
}
