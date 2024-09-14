package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

type Connection struct {
	logger log.Logger
	driver.Conn
}

func New(logger log.Logger, config Config) (*Connection, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: config.Hosts,
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},
		Debug: config.Debug,
		Debugf: func(format string, v ...interface{}) {
			logger.Info().Msgf(format, v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": config.MaxExecutionTime,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      config.DialTimeout,
		MaxOpenConns:     config.MaxOpenConnections,
		MaxIdleConns:     config.MaxIdleConnections,
		ConnMaxLifetime:  config.ConnectionMaxLifetime,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		BlockBufferSize:  uint8(config.BlockBufferSize),
	})
	if err != nil {
		return nil, fmt.Errorf("clickhouse connection failed: %w", err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("clickhouse ping failed %w", err)
	}

	return &Connection{
		logger: logger,
		Conn:   conn,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Conn.Close(); err != nil {
		c.logger.Error().Err(err).Msg("clickhouse failed to close connection")
	}
}
