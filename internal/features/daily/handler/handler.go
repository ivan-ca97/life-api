package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type SummaryHandler interface {
	GetSummary(r *http.Request) (*summaryResponse, int, error)
}

type summaryHandler struct {
	service ports.AuthorizedSummaryService
}

var _ SummaryHandler = (*summaryHandler)(nil)

func NewSummaryHandler(service ports.AuthorizedSummaryService) *summaryHandler {
	return &summaryHandler{
		service: service,
	}
}

func (h *summaryHandler) GetSummary(r *http.Request) (*summaryResponse, int, error) {
	date, err := api.QueryParamDate(r, "date")
	if err != nil {
		return nil, 0, err
	}
	if date == nil {
		return nil, 0, cerr.NewBadRequestError("date query parameter is required (format: YYYY-MM-DD)")
	}
	summary, err := h.service.GetSummary(r.Context(), *date)
	if err != nil {
		return nil, 0, err
	}
	return summaryFromDomain(summary), http.StatusOK, nil
}
