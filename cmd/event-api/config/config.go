package config

import (
	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/http"
	redisqueue "github.com/iamsorryprincess/project-layout/internal/queue/redis"
)

type Config struct {
	LogLevel string

	MySQL mysql.Config

	Redis redis.Config

	Clickhouse clickhouse.Config

	Tarantool tarantool.Config

	HTTP http.Config

	ClicksConsumer redisqueue.ConsumerConfig
}
