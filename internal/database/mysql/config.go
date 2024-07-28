package mysql

import "github.com/iamsorryprincess/project-layout/internal/configuration"

type Config struct {
	ConnectionString      string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime configuration.Duration
	ConnectionMaxIdleTime configuration.Duration
}
