package middleware

import (
	"errors"
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/http/httproute"
)

func Recovery(next httproute.HandlerFunc) httproute.HandlerFunc {
	return func(request *httproute.Request, response *httproute.Response) {
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

				request.LogErrorWithCode(http.StatusInternalServerError, "request recovery from panic: %v", err)
				response.Status(http.StatusInternalServerError)
			}
		}()

		next(request, response)
	}
}
