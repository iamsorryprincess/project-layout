package nats

import "github.com/iamsorryprincess/project-layout/internal/pkg/config"

type Config struct {
	ConnectionString string

	MaxReconnectCount *int
	ReconnectWait     *config.Duration

	Timeout      *config.Duration
	PingInterval *config.Duration
}
