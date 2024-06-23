package clickhouse

import "github.com/iamsorryprincess/project-layout/internal/pkg/config"

type Config struct {
	Hosts                 []string
	User                  string
	Password              string
	Database              string
	Debug                 bool
	MaxExecutionTime      int
	DialTimeout           config.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime config.Duration
	BlockBufferSize       int
}
