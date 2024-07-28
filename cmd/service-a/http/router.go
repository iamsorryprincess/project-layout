package http

import (
	"github.com/iamsorryprincess/project-layout/cmd/service-a/http/handler"
	"github.com/iamsorryprincess/project-layout/cmd/service-a/http/middleware"
	"github.com/iamsorryprincess/project-layout/internal/http/httproute"
	middlewarecommon "github.com/iamsorryprincess/project-layout/internal/http/middleware"
	"github.com/iamsorryprincess/project-layout/internal/log"
)

func NewRouter(dataProvider handler.DataProvider, sessionProvider handler.SessionProvider, logger log.Logger) *httproute.Router {
	router := httproute.NewRouter(logger)

	router.Use(middlewarecommon.Recovery)
	router.Use(middlewarecommon.Cors)
	router.Use(middleware.Test)

	h := handler.NewHandler(logger, sessionProvider, dataProvider)

	router.Route("/test", h.SaveData)

	return router
}
