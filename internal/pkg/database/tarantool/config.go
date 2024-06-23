package tarantool

import "github.com/iamsorryprincess/project-layout/internal/pkg/config"

type Config struct {
	Host              string
	User              string
	Password          string
	Timeout           config.Duration
	ReconnectInterval config.Duration
}
