package http

import (
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/http/request"
	"github.com/iamsorryprincess/project-layout/internal/http/response"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

type HandlerFunc func(request *request.Request, response *response.Response)

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Router struct {
	mux    *http.ServeMux
	logger log.Logger

	middlewares []MiddlewareFunc
}

func NewRouter(logger log.Logger) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.mux.ServeHTTP(writer, request)
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) Route(pattern string, handler HandlerFunc) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	r.mux.HandleFunc(pattern, func(writer http.ResponseWriter, req *http.Request) {
		handler(request.NewRequest(req, r.logger), response.NewResponse(writer, r.logger))
	})
}

func (r *Router) RouteWith(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	r.mux.HandleFunc(pattern, func(writer http.ResponseWriter, req *http.Request) {
		handler(request.NewRequest(req, r.logger), response.NewResponse(writer, r.logger))
	})
}
