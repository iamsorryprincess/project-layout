package middleware

import (
	"errors"
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/log"
)

func Recovery(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			defer func() {
				r := recover()
				if r != nil {
					var err error
					switch t := r.(type) {
					case string:
						err = errors.New(t)
					case error:
						err = t
					default:
						err = errors.New("unknown error")
					}

					logger.Error().
						Str("url", request.RequestURI).
						Str("method", request.Method).
						Int("status", http.StatusInternalServerError).
						Err(err).
						Msg("http request recovered from panic")

					writer.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(writer, request)
		})
	}
}
