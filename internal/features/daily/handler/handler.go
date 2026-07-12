package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type SummaryHandler interface {
	GetSummary(r *http.Request) (*summaryResponse, int, error)
	GetSummaryRange(r *http.Request) (*summaryRangeResponse, int, error)
	GetDailyCheck(r *http.Request) (*dailyCheckResponse, int, error)
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
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := api.QueryParamDate(r, "date")
	if err != nil {
		return nil, 0, err
	}
	if date == nil {
		return nil, 0, cerr.NewBadRequestError("date query parameter is required (format: YYYY-MM-DD)")
	}
	summary, err := h.service.GetSummary(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	return summaryFromDomain(summary), http.StatusOK, nil
}

func (h *summaryHandler) GetDailyCheck(r *http.Request) (*dailyCheckResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := api.QueryParamDate(r, "date")
	if err != nil {
		return nil, 0, err
	}
	if date == nil {
		return nil, 0, cerr.NewBadRequestError("date query parameter is required (format: YYYY-MM-DD)")
	}
	check, err := h.service.GetDailyCheck(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	return dailyCheckFromDomain(check), http.StatusOK, nil
}

func (h *summaryHandler) GetSummaryRange(r *http.Request) (*summaryRangeResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	from, err := api.QueryParamDate(r, "from")
	if err != nil {
		return nil, 0, err
	}
	if from == nil {
		return nil, 0, cerr.NewBadRequestError("from query parameter is required (format: YYYY-MM-DD)")
	}
	to, err := api.QueryParamDate(r, "to")
	if err != nil {
		return nil, 0, err
	}
	if to == nil {
		return nil, 0, cerr.NewBadRequestError("to query parameter is required (format: YYYY-MM-DD)")
	}
	if from.After(*to) {
		return nil, 0, cerr.NewBadRequestError("from must be before or equal to to")
	}
	days := int(to.Sub(*from).Hours()/24) + 1
	if days > 365 {
		return nil, 0, cerr.NewBadRequestError("date range cannot exceed 365 days")
	}
	summaries, err := h.service.GetSummaryRange(r.Context(), userId, *from, *to)
	if err != nil {
		return nil, 0, err
	}
	data := make([]summaryResponse, len(summaries))
	for i, s := range summaries {
		data[i] = *summaryFromDomain(&s)
	}
	result := &summaryRangeResponse{
		Data: data,
	}
	return result, http.StatusOK, nil
}
