package http

import (
	"net/http"

	httputils "github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/http/middleware"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

func NewRouter(logger log.Logger, handler *Handler) http.Handler {
	router := httputils.NewRouter()

	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS)

	router.HandleFunc("/api/click", handler.HandleClick)

	return router
}
