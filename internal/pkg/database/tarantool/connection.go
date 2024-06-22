package tarantool

import (
	"context"
	"fmt"

	"github.com/tarantool/go-tarantool/v2"
)

type Connection struct {
	*tarantool.Connection
}

func New(config Config) (*Connection, error) {
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
		return nil, fmt.Errorf("tarantool connection failed: %v", err)
	}

	if _, err = conn.Do(tarantool.NewPingRequest()).Get(); err != nil {
		return nil, fmt.Errorf("tarantool ping failed: %v", err)
	}

	return &Connection{
		Connection: conn,
	}, nil
}
