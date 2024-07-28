package config

import (
	"github.com/iamsorryprincess/project-layout/internal/configuration"
	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/messaging/nats"
)

type Config struct {
	LogLevel    string
	Timeout     configuration.Duration
	Coefficient float64
	Expire      configuration.Duration
	Interval    configuration.Duration

	HTTP http.Config

	MySQL mysql.Config

	Redis redis.Config

	Clickhouse clickhouse.Config

	Tarantool tarantool.Config

	Nats nats.Config
}
