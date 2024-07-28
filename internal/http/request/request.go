package request

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/rs/zerolog"
)

var ErrJSONContentType = errors.New("invalid Content-Type for json")

type Request struct {
	*http.Request
	logger log.Logger

	IP string
}

func NewRequest(request *http.Request, logger log.Logger) *Request {
	ip := request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = request.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = "127.0.0.1"
		}
	}

	return &Request{
		Request: request,
		logger:  logger,
		IP:      ip,
	}
}

func (r *Request) ParseJSON(result interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrJSONContentType
	}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		r.LogError("failed to close http request body: %v", err)
	}
	return nil
}

func (r *Request) LogDebug(format string, args ...interface{}) {
	r.logEventData(0, format, r.logger.Debug(), args...)
}

func (r *Request) LogDebugWithCode(code int, format string, args ...interface{}) {
	r.logEventData(code, format, r.logger.Debug(), args...)
}

func (r *Request) LogInfo(format string, args ...interface{}) {
	r.logEventData(0, format, r.logger.Info(), args...)
}

func (r *Request) LogInfoWithCode(code int, format string, args ...interface{}) {
	r.logEventData(code, format, r.logger.Info(), args...)
}

func (r *Request) LogWarning(format string, args ...interface{}) {
	r.logEventData(0, format, r.logger.Warn(), args...)
}

func (r *Request) LogWarningWithCode(code int, format string, args ...interface{}) {
	r.logEventData(code, format, r.logger.Warn(), args...)
}

func (r *Request) LogError(format string, args ...interface{}) {
	r.logEventData(0, format, r.logger.Error(), args...)
}

func (r *Request) LogErrorWithCode(code int, format string, args ...interface{}) {
	r.logEventData(code, format, r.logger.Error(), args...)
}

func (r *Request) logEventData(code int, format string, event *zerolog.Event, args ...interface{}) {
	event.Str("type", "http").
		Str("method", r.Method).
		Str("url", r.RequestURI).
		Str("ip", r.IP)

	if code != 0 {
		event.Int("code", code)
	}

	event.Msgf(format, args...)
}
