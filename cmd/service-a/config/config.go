package config

import (
	"github.com/iamsorryprincess/project-layout/internal/pkg/config"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/tarantool"
)

type Config struct {
	LogLevel    string
	Timeout     config.Duration
	Coefficient float64
	Expire      config.Duration

	MySQL mysql.Config

	Redis redis.Config

	Clickhouse clickhouse.Config

	Tarantool tarantool.Config
}

func New(serviceName string) (Config, error) {
	var cfg Config
	if err := config.Parse(serviceName, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
