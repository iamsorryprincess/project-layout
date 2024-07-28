package middleware

import "github.com/iamsorryprincess/project-layout/internal/http/httproute"

func Test(next httproute.HandlerFunc) httproute.HandlerFunc {
	return func(request *httproute.Request, response *httproute.Response) {
		request.LogInfo("test")
		next(request, response)
	}
}
