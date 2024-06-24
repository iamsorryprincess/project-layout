package middleware

import (
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

func Test(logger log.Logger) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			logger.Info().
				Str("type", "http").
				Str("method", request.Method).
				Str("url", request.RequestURI).
				Msg("test")
			handler.ServeHTTP(writer, request)
		})
	}
}
