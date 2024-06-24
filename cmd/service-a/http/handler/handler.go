package handler

import (
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

type DataProvider interface {
	SaveData(data []string) error
}

type Handler struct {
	provider DataProvider
	logger   log.Logger
}

func NewHandler(provider DataProvider, logger log.Logger) *Handler {
	return &Handler{
		provider: provider,
		logger:   logger,
	}
}

func (h *Handler) SaveData(writer http.ResponseWriter, request *http.Request) {
	if err := h.provider.SaveData([]string{request.URL.Query().Get("data")}); err != nil {
		h.logger.Error().
			Str("type", "http").
			Str("method", request.Method).
			Str("url", request.RequestURI).
			Int("code", http.StatusInternalServerError).
			Msgf("save data failed: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
