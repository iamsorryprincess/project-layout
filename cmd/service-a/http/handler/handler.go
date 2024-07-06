package handler

import (
	"context"
	"net/http"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/model"
	"github.com/iamsorryprincess/project-layout/internal/app/domain"
	httputils "github.com/iamsorryprincess/project-layout/internal/pkg/http"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

type SessionProvider interface {
	Get(ctx context.Context, input model.SessionInput) (domain.Session, error)
}

type DataProvider interface {
	SaveData(ctx context.Context, input model.DataInput) error
}

type Handler struct {
	logger          log.Logger
	sessionProvider SessionProvider
	dataProvider    DataProvider
}

func NewHandler(logger log.Logger, sessionProvider SessionProvider, dataProvider DataProvider) *Handler {
	return &Handler{
		logger:          logger,
		sessionProvider: sessionProvider,
		dataProvider:    dataProvider,
	}
}

func (h *Handler) SaveData(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	sessionInput := model.SessionInput{
		Method:      request.Method,
		URL:         request.RequestURI,
		IP:          httputils.ParseIP(request),
		UtmContent:  query.Get("utm_content"),
		UtmTerm:     query.Get("utm_term"),
		UtmCampaign: query.Get("utm_campaign"),
		UtmSource:   query.Get("utm_source"),
		UtmMedium:   query.Get("utm_medium"),
	}

	session, err := h.sessionProvider.Get(request.Context(), sessionInput)
	if err != nil {
		h.logger.Error().
			Str("type", "http").
			Str("method", request.Method).
			Str("url", request.RequestURI).
			Int("code", http.StatusInternalServerError).
			Msgf("get session failed: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	input := model.DataInput{
		Data:    request.URL.Query().Get("data"),
		Session: session,
	}

	if err = h.dataProvider.SaveData(request.Context(), input); err != nil {
		h.logger.Error().
			Str("type", "http").
			Str("method", request.Method).
			Str("url", request.RequestURI).
			Int("code", http.StatusInternalServerError).
			Msgf("save data failed: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
