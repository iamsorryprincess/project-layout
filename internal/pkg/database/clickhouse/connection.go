package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

type Connection struct {
	driver.Conn
}

func New(config Config, logger log.Logger) (*Connection, error) {
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
		return nil, fmt.Errorf("clickhouse connection failed: %v", err)
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("clickhouse ping failed %v", err)
	}

	return &Connection{
		Conn: conn,
	}, nil
}
