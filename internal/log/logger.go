package log

import "github.com/rs/zerolog"

type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event

	With() zerolog.Context
}
