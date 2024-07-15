package config

import (
	"flag"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/pkg/config"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/project-layout/internal/pkg/messaging/nats"
)

type Config struct {
	LogLevel    string
	Timeout     config.Duration
	Coefficient float64
	Expire      config.Duration
	Interval    config.Duration

	HTTP http.Config

	MySQL mysql.Config

	Redis redis.Config

	Clickhouse clickhouse.Config

	Tarantool tarantool.Config

	Nats nats.Config
}

func New(serviceName string) (Config, error) {
	path := flag.String("c", fmt.Sprintf("configs/local/%s.config.json", serviceName), "config path")
	flag.Parse()

	var cfg Config
	if err := config.Parse(*path, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
