package service

import (
	"context"
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

type Repository interface {
	GetData() ([]string, error)
}

type Service struct {
	repository Repository
	logger     log.Logger
}

func NewService(repository Repository, logger log.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) PrintData(_ context.Context) error {
	data, err := s.repository.GetData()
	if err != nil {
		return fmt.Errorf("failed to get data for print: %w", err)
	}
	s.logger.Info().Interface("data", data).Send()
	return nil
}
