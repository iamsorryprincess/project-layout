package http

import "github.com/iamsorryprincess/project-layout/internal/pkg/config"

type Config struct {
	Address string

	ReadTimeout       config.Duration
	ReadHeaderTimeout config.Duration
	WriteTimeout      config.Duration
	IdleTimeout       config.Duration
	ShutdownTimeout   config.Duration

	MaxHeaderBytes int

	DisableGeneralOptionsHandler bool
}
