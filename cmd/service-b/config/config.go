package config

import (
	"github.com/iamsorryprincess/project-layout/internal/configuration"
	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/database/redis"
)

type Config struct {
	LogLevel string

	EventsConsumeCount    int
	EventsConsumeInterval configuration.Duration

	Redis redis.Config

	Clickhouse clickhouse.Config
}
