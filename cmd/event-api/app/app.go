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

	defer a.close()
	defer cancel()

	a.ctx = ctx

	if err := a.initConfig(); err != nil {
		return
	}

	a.logger = log.New(a.config.LogLevel, serviceName)

	if err := a.initDatabases(); err != nil {
		return
	}

	a.httpServer = http.NewServer(a.logger, a.config.HTTP, nil)
	a.httpServer.Start()

	a.logger.Info().Msg("service started")

	s := background.Wait()

	a.logger.Info().Str("stop_signal", s.String()).Msg("service stopped")
}

func (a *App) initConfig() error {
	var err error
	if a.config, err = configuration.New[config.Config](); err != nil {
		log.New("error", serviceName).Error().Err(err).Msg("failed to load configuration")
		return err
	}
	return nil
}

func (a *App) initDatabases() error {
	var err error
	if a.mysqlConn, err = mysql.New(a.logger, a.config.MySQL); err != nil {
		a.logger.Error().Err(err).Msg("failed connect to mysql")
		return err
	}
	a.logger.Info().Msg("mysql successfully connected")

	if a.redisConn, err = redis.New(a.logger, a.config.Redis); err != nil {
		a.logger.Error().Err(err).Msg("failed connect to redis")
		return err
	}
	a.logger.Info().Msg("redis successfully connected")

	if a.clickhouseConn, err = clickhouse.New(a.logger, a.config.Clickhouse); err != nil {
		a.logger.Error().Err(err).Msg("failed connect to clickhouse")
		return err
	}
	a.logger.Info().Msg("clickhouse successfully connected")

	if a.tarantoolConn, err = tarantool.New(a.logger, a.config.Tarantool); err != nil {
		a.logger.Error().Err(err).Msg("failed connect to tarantool")
		return err
	}
	a.logger.Info().Msg("tarantool successfully connected")

	return nil
}

func (a *App) close() {
	if a.httpServer != nil {
		a.httpServer.Stop()
	}

	if a.mysqlConn != nil {
		a.mysqlConn.Close()
	}

	if a.redisConn != nil {
		a.redisConn.Close()
	}

	if a.tarantoolConn != nil {
		a.tarantoolConn.Close()
	}
}
