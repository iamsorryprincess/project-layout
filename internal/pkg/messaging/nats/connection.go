package nats

import (
	"fmt"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/nats-io/nats.go"
)

type Connection struct {
	*nats.Conn
}

func New(config Config, logger log.Logger) (*Connection, error) {
	maxReconnectCount := -1
	if config.MaxReconnectCount != nil {
		maxReconnectCount = *config.MaxReconnectCount
	}

	reconnectWait := time.Second * 2
	if config.ReconnectWait != nil {
		val := *config.ReconnectWait
		reconnectWait = val.Duration
	}

	timeout := time.Second * 2
	if config.Timeout != nil {
		val := *config.Timeout
		timeout = val.Duration
	}

	pingInterval := time.Minute * 2
	if config.PingInterval != nil {
		val := *config.PingInterval
		pingInterval = val.Duration
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
