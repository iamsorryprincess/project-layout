package httproute

import (
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/log"
)

type HandlerFunc func(request *Request, response *Response)

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
		handler(newRequest(req, r.logger), newResponse(writer, r.logger))
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
		handler(newRequest(req, r.logger), newResponse(writer, r.logger))
	})
}
