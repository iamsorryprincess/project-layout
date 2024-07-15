package config

import (
	"flag"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/pkg/config"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
)

type Config struct {
	LogLevel string

	EventsConsumeCount    int
	EventsConsumeInterval config.Duration

	Redis redis.Config

	Clickhouse clickhouse.Config
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
