package http

import (
	"errors"
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

func Recovery(logger log.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
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
						Str("type", "http").
						Str("url", request.RequestURI).
						Int("code", http.StatusInternalServerError).
						Msgf("request recovery from panic: %v", err)

					writer.WriteHeader(http.StatusInternalServerError)
				}
			}()

			handler.ServeHTTP(writer, request)
		})
	}
}
