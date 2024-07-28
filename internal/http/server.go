package http

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/iamsorryprincess/project-layout/internal/log"
)

type Server struct {
	config Config
	logger log.Logger
	server *http.Server
	wg     sync.WaitGroup
}

func NewServer(config Config, logger log.Logger, handler http.Handler) *Server {
	return &Server{
		config: config,
		logger: logger,
		server: &http.Server{
			Addr:                         config.Address,
			Handler:                      handler,
			ReadTimeout:                  config.ReadTimeout.Duration,
			ReadHeaderTimeout:            config.ReadHeaderTimeout.Duration,
			WriteTimeout:                 config.WriteTimeout.Duration,
			IdleTimeout:                  config.IdleTimeout.Duration,
			MaxHeaderBytes:               config.MaxHeaderBytes,
			DisableGeneralOptionsHandler: config.DisableGeneralOptionsHandler,
		},
		wg: sync.WaitGroup{},
	}
}

func (s *Server) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal().Str("type", "http").Msgf("http server listen error: %v", err)
		}
	}()
}

func (s *Server) Stop() {
	ctx := context.Background()

	if s.config.ShutdownTimeout.Duration > 0 {
		newCtx, cancel := context.WithTimeout(ctx, s.config.ShutdownTimeout.Duration)
		defer cancel()
		ctx = newCtx
	}

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error().Str("type", "http").Msgf("http server shutdown error: %v", err)
	}

	s.wg.Wait()
}
