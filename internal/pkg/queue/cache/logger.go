package cache

import (
	"fmt"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

func newLogger(key string, subtype string, logger log.Logger) log.Logger {
	loggerWithFields := logger.With().
		Str("type", fmt.Sprintf("file_queue_%s", subtype)).
		Str("key", key).
		Logger()
	return log.Logger{
		Logger: loggerWithFields,
	}
}
