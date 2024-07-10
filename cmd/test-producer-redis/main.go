package main

import (
	"net/http"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/app/domain"
	"github.com/iamsorryprincess/project-layout/internal/pkg/background"
	"github.com/iamsorryprincess/project-layout/internal/pkg/database/redis"
	httputils "github.com/iamsorryprincess/project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	redisqueue "github.com/iamsorryprincess/project-layout/internal/pkg/queue/redis"
)

const serviceName = "test-producer-redis"

func main() {
	logger := log.New("info", serviceName)

	redisConn, err := redis.New(redis.Config{Host: "localhost:6379"}, logger)
	if err != nil {
		logger.Error().Msg(err.Error())
		return
	}

	defer redisConn.Close()

	producer := redisqueue.NewProducer[domain.Event]("events", redisConn)

	router := http.NewServeMux()

	router.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		event := domain.Event{
			CreatedAt:  time.Now(),
			IP:         httputils.ParseIP(request),
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

	server := httputils.NewServer(httputils.Config{Address: ":8080"}, logger, router)
	server.Start()
	defer server.Stop()

	logger.Info().Msg("service started")
	background.Wait()
	logger.Info().Msg("service stopped")
}
