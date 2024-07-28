package clickhouse

import "github.com/iamsorryprincess/project-layout/internal/configuration"

type Config struct {
	Hosts                 []string
	User                  string
	Password              string
	Database              string
	Debug                 bool
	MaxExecutionTime      int
	DialTimeout           configuration.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime configuration.Duration
	BlockBufferSize       int
}
