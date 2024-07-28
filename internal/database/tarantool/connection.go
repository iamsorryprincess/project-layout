package tarantool

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/tarantool/go-tarantool/v2"
)

type Connection struct {
	*tarantool.Connection

	logger log.Logger
}

func New(config Config, logger log.Logger) (*Connection, error) {
	dialer := tarantool.NetDialer{
		Address:  config.Host,
		User:     config.User,
		Password: config.Password,
	}

	conn, err := tarantool.Connect(context.Background(), dialer, tarantool.Opts{
		Timeout:   config.Timeout.Duration,
		Reconnect: config.ReconnectInterval.Duration,
	})
	if err != nil {
		return nil, fmt.Errorf("tarantool connection failed: %v", err)
	}

	if _, err = conn.Do(tarantool.NewPingRequest()).Get(); err != nil {
		return nil, fmt.Errorf("tarantool ping failed: %v", err)
	}

	return &Connection{
		Connection: conn,
		logger:     logger,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Connection.Close(); err != nil {
		c.logger.Error().Str("type", "tarantool").Msg("tarantool failed to close connection")
	}
}
