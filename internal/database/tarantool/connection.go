package tarantool

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/tarantool/go-tarantool/v2"
)

type Connection struct {
	logger log.Logger
	*tarantool.Connection
}

func New(logger log.Logger, config Config) (*Connection, error) {
	dialer := tarantool.NetDialer{
		Address:  config.Host,
		User:     config.User,
		Password: config.Password,
	}

	conn, err := tarantool.Connect(context.Background(), dialer, tarantool.Opts{
		Timeout:   config.Timeout,
		Reconnect: config.ReconnectInterval,
	})
	if err != nil {
		return nil, fmt.Errorf("tarantool connection failed: %w", err)
	}

	if _, err = conn.Do(tarantool.NewPingRequest()).Get(); err != nil {
		return nil, fmt.Errorf("tarantool ping failed: %w", err)
	}

	return &Connection{
		Connection: conn,
		logger:     logger,
	}, nil
}

func (c *Connection) Close() {
	if err := c.Connection.Close(); err != nil {
		c.logger.Error().Err(err).Msg("tarantool failed to close connection")
	}
}
