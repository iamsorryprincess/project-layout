package main

import (
	"context"
	"net/http"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/background"
	"github.com/iamsorryprincess/project-layout/internal/database/redis"
	"github.com/iamsorryprincess/project-layout/internal/domain"
	httputils "github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/iamsorryprincess/project-layout/internal/queue"
	redisqueue "github.com/iamsorryprincess/project-layout/internal/queue/redis"
)

const serviceName = "test-producer-cache"

func main() {
	logger := log.New("info", serviceName)

	redisConn, err := redis.New(redis.Config{Host: "localhost:6379"}, logger)
	if err != nil {
		logger.Error().Msg(err.Error())
		return
	}

	defer redisConn.Close()

	redisProducer := redisqueue.NewProducer[domain.Event]("events", redisConn)
	producer := queue.NewFileCachingProducer[domain.Event]("events", logger, redisProducer)

	router := http.NewServeMux()

	router.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		event := domain.Event{
			CreatedAt:  time.Now(),
			IP:         "",
			CountryID:  "RU",
			PlatformID: 1,
		}

		if pErr := producer.Produce(request.Context(), event); pErr != nil {
			logger.Error().Msg(pErr.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	worker := background.NewWorker(logger)

	if _, err = worker.StartWithInterval(context.Background(), "send events", time.Second*3, producer.Send); err != nil {
		logger.Error().Msg(err.Error())
		return
	}

	server := httputils.NewServer(httputils.Config{Address: ":8081"}, logger, router)
	server.Start()

	logger.Info().Msg("service started")
	background.Wait()

	server.Stop()

	if err = producer.Send(context.Background()); err != nil {
		logger.Error().Msg(err.Error())
	}

	worker.StopAll()

	logger.Info().Msg("service stopped")
}
