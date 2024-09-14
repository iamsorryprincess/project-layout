package http

import "net/http"

type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	if len(r.middlewares) == 0 {
		r.mux.HandleFunc(pattern, handler)
		return
	}

	var h http.Handler
	h = handler

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	r.mux.Handle(pattern, h)
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.mux.ServeHTTP(writer, request)
}
