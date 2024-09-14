package tarantool

import "time"

type Config struct {
	Host              string
	User              string
	Password          string
	Timeout           time.Duration
	ReconnectInterval time.Duration
}
