package app

import (
	"context"
	"time"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/config"
	httptransport "github.com/iamsorryprincess/project-layout/cmd/service-a/http"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/repository"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/service"
	"github.com/iamsorryprincess/project-layout/internal/app/domain"
	sessionrepository "github.com/iamsorryprincess/project-layout/internal/app/session/repository"
	"github.com/iamsorryprincess/project-layout/internal/pkg/background"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/mysql"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/tarantool"
	"github.com/iamsorryprincess/project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/messaging/nats"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue/cache"
	redisqueue "github.com/iamsorryprincess/project-layout/internal/pkg/queue/redis"
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

	dataRepository    *repository.Repository
	sessionRepository *sessionrepository.Repository

	sessionProducer queue.Producer[domain.Session]
	eventProducer   *cache.Producer[domain.Event]

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

func (a *App) initNats() {
	var err error
	if a.natsConn, err = nats.New(a.config.Nats, a.logger); err != nil {
		a.logger.Fatal().Str("type", "nats").Msg(err.Error())
	}
	a.logger.Info().Str("type", "nats").Msg("nats connected")
}

func (a *App) initRepositories() {
	a.dataRepository = repository.New()
	a.sessionRepository = sessionrepository.NewRepository("session", time.Minute*15, a.redisConn)
}

func (a *App) initQueue() {
	a.sessionProducer = redisqueue.NewProducer[domain.Session]("sessions", a.redisConn)
	a.eventProducer = cache.NewProducer[domain.Event]("events", a.logger, redisqueue.NewProducer[domain.Event]("events", a.redisConn))
}

func (a *App) initServices() {
	a.sessionService = service.NewSessionService(a.logger, a.sessionRepository, a.sessionProducer)
	a.dataService = service.NewDataService(a.logger, a.eventProducer, a.dataRepository)
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger)
	if _, err := a.worker.StartWithInterval(a.ctx, "sending event messages", time.Minute, a.eventProducer.Send); err != nil {
		a.logger.Fatal().Str("type", "worker").Msg("failed to start sending event messages worker")
	}
}

func (a *App) initHTTP() {
	router := httptransport.NewRouter(a.dataService, a.sessionService, a.logger)
	a.httpServer = http.NewServer(a.config.HTTP, a.logger, router)
	a.httpServer.Start()
}

func (a *App) close() {
	a.httpServer.Stop()

	if err := a.eventProducer.Send(context.Background()); err != nil {
		a.logger.Error().Str("type", "events_producer").Msgf("failed to send events: %v", err)
	}

	a.worker.StopAll()
	a.natsConn.Close()
	a.mysqlConn.Close()
	a.redisConn.Close()
	a.clickhouseConn.Close()
	a.tarantoolConn.Close()
}
