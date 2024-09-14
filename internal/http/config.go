package http

import "time"

type Config struct {
	Address string

	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration

	MaxHeaderBytes int

	DisableGeneralOptionsHandler bool
}
