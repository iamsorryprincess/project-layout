package nats

import "github.com/iamsorryprincess/project-layout/internal/configuration"

type Config struct {
	ConnectionString string

	MaxReconnectCount int
	ReconnectWait     configuration.Duration

	Timeout      configuration.Duration
	PingInterval configuration.Duration
}
