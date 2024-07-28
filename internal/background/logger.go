package background

import "github.com/iamsorryprincess/project-layout/internal/log"

func newLogger(logger log.Logger) log.Logger {
	loggerWithFields := logger.With().
		Str("type", "worker").
		Logger()
	return log.Logger{
		Logger: loggerWithFields,
	}
}
