package http

import "github.com/iamsorryprincess/project-layout/internal/configuration"

type Config struct {
	Address string

	ReadTimeout       configuration.Duration
	ReadHeaderTimeout configuration.Duration
	WriteTimeout      configuration.Duration
	IdleTimeout       configuration.Duration
	ShutdownTimeout   configuration.Duration

	MaxHeaderBytes int

	DisableGeneralOptionsHandler bool
}
