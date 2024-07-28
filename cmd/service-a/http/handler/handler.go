package handler

import (
	"context"
	"net/http"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/model"
	"github.com/iamsorryprincess/project-layout/internal/domain"
	"github.com/iamsorryprincess/project-layout/internal/http/request"
	"github.com/iamsorryprincess/project-layout/internal/http/response"
	"github.com/iamsorryprincess/project-layout/internal/log"
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

func (h *Handler) SaveData(request *request.Request, response *response.Response) {
	query := request.URL.Query()
	sessionInput := model.SessionInput{
		Method:      request.Method,
		URL:         request.RequestURI,
		IP:          request.IP,
		UtmContent:  query.Get("utm_content"),
		UtmTerm:     query.Get("utm_term"),
		UtmCampaign: query.Get("utm_campaign"),
		UtmSource:   query.Get("utm_source"),
		UtmMedium:   query.Get("utm_medium"),
	}

	session, err := h.sessionProvider.Get(request.Context(), sessionInput)
	if err != nil {
		request.LogErrorWithCode(http.StatusInternalServerError, "get session failed: %v", err)
		response.Status(http.StatusInternalServerError)
	}

	input := model.DataInput{
		Data:    request.URL.Query().Get("data"),
		Session: session,
	}

	if err = h.dataProvider.SaveData(request.Context(), input); err != nil {
		request.LogErrorWithCode(http.StatusInternalServerError, "save data failed: %v", err)
		response.Status(http.StatusInternalServerError)
		return
	}

	response.Status(http.StatusOK)
}
