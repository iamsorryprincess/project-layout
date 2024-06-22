package clickhouse

import "time"

type Config struct {
	Hosts                 []string
	User                  string
	Password              string
	Database              string
	Debug                 bool
	MaxExecutionTime      int
	DialTimeout           time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	BlockBufferSize       int
}
