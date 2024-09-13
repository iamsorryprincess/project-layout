package app

import (
	"context"

	"github.com/iamsorryprincess/project-layout/cmd/event-api/config"
	"github.com/iamsorryprincess/project-layout/internal/background"
	"github.com/iamsorryprincess/project-layout/internal/configuration"
	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

const serviceName = "event-api"

type App struct {
	ctx    context.Context
	logger log.Logger

	config config.Config

	mysqlConn      *mysql.Connection
	redisConn      *redis.Connection
	clickhouseConn *clickhouse.Connection
	tarantoolConn  *tarantool.Connection

	httpServer *http.Server
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.ctx = ctx

	var err error
	if a.config, err = configuration.New[config.Config](); err != nil {
		log.New("error", serviceName).Error().Err(err).Msg("failed to load configuration")
		return
	}

	a.logger = log.New(a.config.LogLevel, serviceName)

	a.httpServer = http.NewServer(a.logger, a.config.HTTP, nil)
	a.httpServer.Start()

	a.logger.Info().Msg("service started")

	s := background.Wait()

	a.httpServer.Stop()

	a.logger.Info().Str("stop_signal", s.String()).Msg("service stopped")
}
