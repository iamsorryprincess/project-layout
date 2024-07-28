package nats

import (
	"fmt"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/nats-io/nats.go"
)

type Connection struct {
	*nats.Conn
}

func New(config Config, logger log.Logger) (*Connection, error) {
	maxReconnectCount := -1
	if config.MaxReconnectCount != 0 {
		maxReconnectCount = config.MaxReconnectCount
	}

	reconnectWait := time.Second * 2
	if config.ReconnectWait.Duration != 0 {
		reconnectWait = config.ReconnectWait.Duration
	}

	timeout := time.Second * 2
	if config.Timeout.Duration != 0 {
		timeout = config.Timeout.Duration
	}

	pingInterval := time.Minute * 2
	if config.PingInterval.Duration != 0 {
		pingInterval = config.PingInterval.Duration
	}

	conn, err := nats.Connect(config.ConnectionString,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(maxReconnectCount),
		nats.ReconnectWait(reconnectWait),
		nats.Timeout(timeout),
		nats.PingInterval(pingInterval),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Error().Str("type", "nats").Msgf("nats connect failed; trying to reconnect")
		}))
	if err != nil {
		return nil, fmt.Errorf("nats connection failed: %w", err)
	}

	return &Connection{
		Conn: conn,
	}, nil
}
