package middleware

import (
	"errors"
	httpnet "net/http"

	"github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/http/request"
	"github.com/iamsorryprincess/project-layout/internal/http/response"
)

func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(request *request.Request, response *response.Response) {
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

				request.LogErrorWithCode(httpnet.StatusInternalServerError, "request recovery from panic: %v", err)
				response.Status(httpnet.StatusInternalServerError)
			}
		}()

		next(request, response)
	}
}
