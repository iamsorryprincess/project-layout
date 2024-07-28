package service

import (
	"context"
	"time"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/model"
	"github.com/iamsorryprincess/project-layout/internal/domain"
	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/iamsorryprincess/project-layout/internal/queue"
)

type DataService struct {
	logger   log.Logger
	producer queue.Producer[domain.Event]
}

func NewDataService(logger log.Logger, producer queue.Producer[domain.Event]) *DataService {
	return &DataService{
		logger:   logger,
		producer: producer,
	}
}

func (s *DataService) SaveData(ctx context.Context, input model.DataInput) error {
	events := make([]domain.Event, 100000)
	now := time.Now()
	for i := range events {
		events[i] = domain.Event{
			CreatedAt:  now,
			IP:         input.Session.IP,
			CountryID:  input.Session.CountryID,
			PlatformID: input.Session.PlatformID,
		}
	}

	if err := s.producer.Produce(ctx, events...); err != nil {
		s.logger.Error().
			Str("type", "data_service").
			Str("ip", input.Session.IP).
			Str("country", input.Session.CountryID).
			Uint8("platform", input.Session.PlatformID).
			Msgf("failed to produce events: %v", err)
	}

	return nil
}
