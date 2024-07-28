package app

import (
	"context"
	"time"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/config"
	httptransport "github.com/iamsorryprincess/project-layout/cmd/service-a/http"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/service"
	"github.com/iamsorryprincess/project-layout/internal/background"
	"github.com/iamsorryprincess/project-layout/internal/configuration"
	"github.com/iamsorryprincess/project-layout/internal/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/domain"
	"github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/iamsorryprincess/project-layout/internal/messaging/nats"
	"github.com/iamsorryprincess/project-layout/internal/queue"
	redisqueue "github.com/iamsorryprincess/project-layout/internal/queue/redis"
	sessionrepository "github.com/iamsorryprincess/project-layout/internal/session/repository"
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

	natsConn *nats.Connection

	sessionRepository *sessionrepository.Repository

	sessionProducer queue.Producer[domain.Session]
	eventProducer   queue.Producer[domain.Event]

	sessionService *service.SessionService
	dataService    *service.DataService

	worker *background.Worker

	httpServer *http.Server
}

func New() *App {
	return &App{}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx

	a.initConfig()

	a.logger = log.New(a.config.LogLevel, serviceName)

	a.initDatabases()

	a.initNats()

	a.initRepositories()

	a.initQueue()

	a.initServices()

	a.initWorkers()

	a.initHTTP()

	a.logger.Info().Interface("configuration", a.config).Msg("service started")

	stopSignal := background.Wait()

	cancel()
	a.close()

	a.logger.Info().Str("stop_signal", stopSignal.String()).Msg("service stopped")
}

func (a *App) initConfig() {
	if err := configuration.Parse(configuration.TypeJSON, serviceName, &a.config); err != nil {
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

func (a *App) initNats() {
	var err error
	if a.natsConn, err = nats.New(a.config.Nats, a.logger); err != nil {
		a.logger.Fatal().Str("type", "nats").Msg(err.Error())
	}
	a.logger.Info().Str("type", "nats").Msg("nats connected")
}

func (a *App) initRepositories() {
	a.sessionRepository = sessionrepository.NewRepository("session", time.Minute*15, a.redisConn)
}

func (a *App) initQueue() {
	a.sessionProducer = redisqueue.NewProducer[domain.Session]("sessions", a.redisConn)
	a.eventProducer = redisqueue.NewProducer[domain.Event]("events", a.redisConn)
}

func (a *App) initServices() {
	a.sessionService = service.NewSessionService(a.logger, a.sessionRepository, a.sessionProducer)
	a.dataService = service.NewDataService(a.logger, a.eventProducer)
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger)
}

func (a *App) initHTTP() {
	router := httptransport.NewRouter(a.dataService, a.sessionService, a.logger)
	a.httpServer = http.NewServer(a.config.HTTP, a.logger, router)
	a.httpServer.Start()
}

func (a *App) close() {
	a.httpServer.Stop()
	a.worker.StopAll()
	a.natsConn.Close()
	a.mysqlConn.Close()
	a.redisConn.Close()
	a.clickhouseConn.Close()
	a.tarantoolConn.Close()
}
