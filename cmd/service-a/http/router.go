package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/http/handler"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/http/middleware"
	httppkg "github.com/iamsorryprincess/project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

func NewRouter(dataProvider handler.DataProvider, sessionProvider handler.SessionProvider, logger log.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(httppkg.Recovery(logger))
	router.Use(httppkg.Cors)
	router.Use(middleware.Test(logger))

	h := handler.NewHandler(logger, sessionProvider, dataProvider)

	router.Get("/test", h.SaveData)

	return router
}
