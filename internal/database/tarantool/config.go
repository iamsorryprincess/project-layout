package tarantool

import "github.com/iamsorryprincess/project-layout/internal/configuration"

type Config struct {
	Host              string
	User              string
	Password          string
	Timeout           configuration.Duration
	ReconnectInterval configuration.Duration
}
