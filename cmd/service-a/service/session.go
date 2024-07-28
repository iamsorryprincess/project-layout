package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/model"
	"github.com/iamsorryprincess/project-layout/internal/domain"
	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/iamsorryprincess/project-layout/internal/queue"
)

type SessionProvider interface {
	Set(ctx context.Context, session domain.Session) error
	Get(ctx context.Context, ip string) (domain.Session, error)
	UpdateTTL(ctx context.Context, ip string) error
}

type SessionService struct {
	logger   log.Logger
	provider SessionProvider
	producer queue.Producer[domain.Session]
}

func NewSessionService(logger log.Logger, provider SessionProvider, producer queue.Producer[domain.Session]) *SessionService {
	return &SessionService{
		logger:   logger,
		provider: provider,
		producer: producer,
	}
}

func (s *SessionService) Get(ctx context.Context, input model.SessionInput) (domain.Session, error) {
	session, err := s.provider.Get(ctx, input.IP)
	if err != nil {
		notFound := domain.ErrNotFound{}
		if !errors.As(err, &notFound) {
			return domain.Session{}, err
		}

		session = domain.Session{
			ID:       uuid.New().String(),
			UserID:   "test",
			Datetime: time.Now(),

			IP:         input.IP,
			CountryID:  "RU",
			PlatformID: 1,

			UtmContent:  input.UtmContent,
			UtmTerm:     input.UtmTerm,
			UtmCampaign: input.UtmCampaign,
			UtmSource:   input.UtmSource,
			UtmMedium:   input.UtmMedium,
		}

		if err = s.producer.Produce(ctx, session); err != nil {
			s.logger.Error().
				Str("type", "http").
				Str("method", input.Method).
				Str("url", input.URL).
				Msgf("produce session failed: %v", err)
		}

		if err = s.provider.Set(ctx, session); err != nil {
			return domain.Session{}, err
		}

		return session, nil
	}

	if err = s.provider.UpdateTTL(ctx, input.IP); err != nil {
		s.logger.Error().
			Str("type", "http").
			Str("method", input.Method).
			Str("url", input.URL).
			Msgf("update session ttl failed: %v", err)
	}

	return session, nil
}
