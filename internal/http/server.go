package http

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/iamsorryprincess/project-layout/internal/log"
)

type Server struct {
	logger log.Logger
	config Config
	server *http.Server

	isRunning bool

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewServer(logger log.Logger, config Config, handler http.Handler) *Server {
	return &Server{
		config: config,
		logger: logger,
		server: &http.Server{
			Addr:                         config.Address,
			Handler:                      handler,
			ReadTimeout:                  config.ReadTimeout,
			ReadHeaderTimeout:            config.ReadHeaderTimeout,
			WriteTimeout:                 config.WriteTimeout,
			IdleTimeout:                  config.IdleTimeout,
			MaxHeaderBytes:               config.MaxHeaderBytes,
			DisableGeneralOptionsHandler: config.DisableGeneralOptionsHandler,
		},

		isRunning: false,

		mu: sync.Mutex{},
		wg: sync.WaitGroup{},
	}
}

func (s *Server) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		s.logger.Warn().Msg("http server is already running")
		return
	}

	s.wg.Add(1)
	s.isRunning = true

	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				s.logger.Info().Msg("http server successfully shutdown")
				return
			}
			s.logger.Error().Err(err).Msgf("http server listen error")
		}
	}()
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	s.isRunning = false

	ctx := context.Background()
	if s.config.ShutdownTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.ShutdownTimeout)
		defer cancel()
	}

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("http server shutdown error")
	}

	s.wg.Wait()
}
