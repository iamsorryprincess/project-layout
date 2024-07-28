package middleware

import (
	"github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/http/request"
	"github.com/iamsorryprincess/project-layout/internal/http/response"
)

func Test(next http.HandlerFunc) http.HandlerFunc {
	return func(request *request.Request, response *response.Response) {
		request.LogInfo("test")
		next(request, response)
	}
}
