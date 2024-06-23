package app

import (
	"context"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/config"
	"github.com/iamsorryprincess/project-layout/internal/pkg/background"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

const serviceName = "service-a"

type App struct {
	ctx    context.Context
	config config.Config
	logger log.Logger

	mysqlConn      *mysql.Connection
	redisConn      *redis.Connection
	clickhouseConn *clickhouse.Connection
	tarantoolConn  *tarantool.Connection
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	a.ctx = context.Background()

	a.initConfig()

	a.logger = log.New(a.config.LogLevel, serviceName)

	a.initDatabases()

	a.logger.Info().Interface("configuration", a.config).Msg("service started")

	defer a.close()

	background.Wait(a.logger)
}

func (a *App) initConfig() {
	var err error
	a.config, err = config.New(serviceName)

	if err != nil {
		logger := log.New("fatal", serviceName)
		logger.Fatal().Str("type", "config").Msgf("failed to load config: %v", err)
	}
}

func (a *App) initDatabases() {
	var err error
	if a.mysqlConn, err = mysql.New(a.config.MySQL, a.logger); err != nil {
		a.logger.Fatal().Str("type", "mysql").Msg(err.Error())
	}

	a.logger.Info().Msg("mysql connected")

	if a.redisConn, err = redis.New(a.config.Redis, a.logger); err != nil {
		a.logger.Fatal().Str("type", "redis").Msg(err.Error())
	}

	a.logger.Info().Msg("redis connected")

	if a.clickhouseConn, err = clickhouse.New(a.config.Clickhouse, a.logger); err != nil {
		a.logger.Fatal().Str("type", "clickhouse").Msg(err.Error())
	}

	a.logger.Info().Msg("clickhouse connected")

	if a.tarantoolConn, err = tarantool.New(a.config.Tarantool, a.logger); err != nil {
		a.logger.Fatal().Str("type", "tarantool").Msg(err.Error())
	}

	a.logger.Info().Msg("tarantool connected")
}

func (a *App) close() {
	a.mysqlConn.Close()
	a.redisConn.Close()
	a.clickhouseConn.Close()
	a.tarantoolConn.Close()
}
