package service

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/model"
	"github.com/iamsorryprincess/project-layout/internal/app/domain"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
	"github.com/iamsorryprincess/project-layout/internal/pkg/queue"
)

type Repository interface {
	SaveData(data string) error
	GetData() ([]string, error)
}

type DataService struct {
	logger     log.Logger
	producer   queue.Producer[domain.Event]
	repository Repository
}

func NewDataService(logger log.Logger, producer queue.Producer[domain.Event], repository Repository) *DataService {
	return &DataService{
		logger:     logger,
		producer:   producer,
		repository: repository,
	}
}

func (s *DataService) SaveData(ctx context.Context, input model.DataInput) error {
	events := make([]domain.Event, 100)
	for i := range events {
		events[i] = domain.Event{
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
			Int("platform", input.Session.PlatformID).
			Msgf("failed to produce events: %v", err)
	}

	if err := s.repository.SaveData(input.Data); err != nil {
		return fmt.Errorf("failed to save data: %w", err)
	}

	return nil
}

func (s *DataService) PrintData(_ context.Context) error {
	data, err := s.repository.GetData()
	if err != nil {
		return fmt.Errorf("failed to get data for print: %w", err)
	}
	s.logger.Info().Interface("data", data).Send()
	return nil
}
