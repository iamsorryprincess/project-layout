package http

import (
	"net/http"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/domain"
	"github.com/iamsorryprincess/project-layout/internal/log"
	"github.com/iamsorryprincess/project-layout/internal/queue"
)

type Handler struct {
	logger log.Logger

	clickProducer queue.Producer[domain.Click]
}

func NewHandler(logger log.Logger, clickProducer queue.Producer[domain.Click]) *Handler {
	return &Handler{
		logger:        logger,
		clickProducer: clickProducer,
	}
}

func (h *Handler) HandleClick(writer http.ResponseWriter, request *http.Request) {
	click := domain.Click{
		CreatedAt: time.Now(),
		IP:        request.Header.Get("X-Real-IP"),
	}

	if err := h.clickProducer.Produce(request.Context(), click); err != nil {
		h.logger.Error().
			Str("url", request.RequestURI).
			Str("method", request.Method).
			Int("status", http.StatusInternalServerError).
			Err(err).
			Msg("failed to produce click")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
