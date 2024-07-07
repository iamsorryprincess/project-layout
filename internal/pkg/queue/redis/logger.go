package redis

import "github.com/iamsorryprincess/project-layout/internal/pkg/log"

func newLogger(key string, logger log.Logger) log.Logger {
	loggerWithFields := logger.With().
		Str("type", "redis_queue").
		Str("key", key).
		Logger()
	return log.Logger{
		Logger: loggerWithFields,
	}
}
