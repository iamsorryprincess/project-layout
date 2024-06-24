package app

import (
	"context"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/config"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/repository"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/service"
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

	repository *repository.Repository

	service *service.Service

	worker *background.Worker
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	a.ctx = context.Background()

	a.initConfig()

	a.logger = log.New(a.config.LogLevel, serviceName)

	a.initDatabases()

	a.initRepositories()

	a.initServices()

	a.initWorkers()

	a.logger.Info().Interface("configuration", a.config).Msg("service started")

	stopSignal := background.Wait()

	a.close()

	a.logger.Info().Str("stop_signal", stopSignal.String()).Msg("service stopped")
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

	a.logger.Info().Str("type", "mysql").Msg("mysql connected")

	if a.redisConn, err = redis.New(a.config.Redis, a.logger); err != nil {
		a.logger.Fatal().Str("type", "redis").Msg(err.Error())
	}

	a.logger.Info().Str("type", "redis").Msg("redis connected")

	if a.clickhouseConn, err = clickhouse.New(a.config.Clickhouse, a.logger); err != nil {
		a.logger.Fatal().Str("type", "clickhouse").Msg(err.Error())
	}

	a.logger.Info().Str("type", "clickhouse").Msg("clickhouse connected")

	if a.tarantoolConn, err = tarantool.New(a.config.Tarantool, a.logger); err != nil {
		a.logger.Fatal().Str("type", "tarantool").Msg(err.Error())
	}

	a.logger.Info().Str("type", "tarantool").Msg("tarantool connected")
}

func (a *App) initRepositories() {
	a.repository = repository.New()
}

func (a *App) initServices() {
	a.service = service.NewService(a.repository, a.logger)
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger)
	if _, err := a.worker.StartWithInterval(a.ctx, "printing data", a.config.Interval.Duration, a.service.PrintData); err != nil {
		a.logger.Fatal().Str("type", "worker").Msg("failed to start printing data worker")
	}
}

func (a *App) close() {
	a.worker.StopAll()
	a.mysqlConn.Close()
	a.redisConn.Close()
	a.clickhouseConn.Close()
	a.tarantoolConn.Close()
}
