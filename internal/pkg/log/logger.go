package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// projectName Set actual project name
const projectName = "project_name"

type Logger struct {
	zerolog.Logger
}

func New(level string, serviceName string) Logger {
	zerolog.TimestampFieldName = "datetime"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldName = "level"
	zerolog.LevelDebugValue = "debug"
	zerolog.LevelInfoValue = "info"
	zerolog.LevelWarnValue = "warning"
	zerolog.LevelErrorValue = "error"
	zerolog.MessageFieldName = "description"

	logLevel := zerolog.InfoLevel
	isNotParsed := false

	switch level {
	case "trace":
		logLevel = zerolog.TraceLevel
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn", "warning":
		logLevel = zerolog.WarnLevel
	case "error", "err":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	case "disable", "disabled":
		logLevel = zerolog.Disabled
	default:
		isNotParsed = true
	}

	logger := zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Str("project", projectName).
		Str("service", serviceName).
		Logger()

	if isNotParsed {
		logger.Warn().Msgf("unknown log level: %s; creating logger with info level as default value", level)
	}

	return Logger{
		Logger: logger,
	}
}
