package mysql

import "github.com/iamsorryprincess/project-layout/internal/pkg/config"

type Config struct {
	ConnectionString      string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime config.Duration
	ConnectionMaxIdleTime config.Duration
}
