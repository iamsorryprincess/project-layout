package app

import (
	"context"

	"github.com/iamsorryprincess/project-layout/cmd/service-b/config"
	"github.com/iamsorryprincess/project-layout/cmd/service-b/repository"
	"github.com/iamsorryprincess/project-layout/internal/app/domain"
	"github.com/iamsorryprincess/project-layout/internal/pkg/background"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/clickhouse"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue/cache"
	redisqueue "github.com/iamsorryprincess/project-layout/internal/pkg/queue/redis"
)

const serviceName = "service-b"

type App struct {
	ctx    context.Context
	config config.Config
	logger log.Logger

	redisConn      *redis.Connection
	clickhouseConn *clickhouse.Connection

	eventRepository *repository.EventRepository

	eventConsumer *cache.Consumer[domain.Event]

	worker *background.Worker
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

	a.initRepositories()

	a.initQueue()

	a.initWorkers()

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
	if a.redisConn, err = redis.New(a.config.Redis, a.logger); err != nil {
		a.logger.Fatal().Str("type", "redis").Msg(err.Error())
	}

	a.logger.Info().Str("type", "redis").Msg("redis connected")

	if a.clickhouseConn, err = clickhouse.New(a.config.Clickhouse, a.logger); err != nil {
		a.logger.Fatal().Str("type", "clickhouse").Msg(err.Error())
	}

	a.logger.Info().Str("type", "clickhouse").Msg("clickhouse connected")
}

func (a *App) initRepositories() {
	a.eventRepository = repository.NewEventRepository(a.clickhouseConn)
}

func (a *App) initQueue() {
	eventProducer := redisqueue.NewProducer[domain.Event]("events", a.redisConn)
	redisEventConsumer := redisqueue.NewConsumer[domain.Event]("events", a.config.EventsConsumeCount, a.logger, a.redisConn)
	a.eventConsumer = cache.NewConsumer[domain.Event]("events", a.config.EventsConsumeCount, a.logger, a.eventRepository, eventProducer, redisEventConsumer)
}

func (a *App) initWorkers() {
	a.worker = background.NewWorker(a.logger)
	if _, err := a.worker.StartWithInterval(a.ctx, "consuming events", a.config.EventsConsumeInterval.Duration, a.eventConsumer.Consume); err != nil {
		a.logger.Fatal().Str("type", "worker").Msg("failed to start consuming events worker")
	}
}

func (a *App) close() {
	a.worker.StopAll()
	a.redisConn.Close()
	a.clickhouseConn.Close()
}
