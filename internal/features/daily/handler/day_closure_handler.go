package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type dayClosureResponse struct {
	Date   string `json:"date"`
	Closed bool   `json:"closed"`
}

type DayClosureHandler interface {
	Close(r *http.Request) (*dayClosureResponse, int, error)
	Open(r *http.Request) (*api.NoResponse, int, error)
	GetStatus(r *http.Request) (*dayClosureResponse, int, error)
}

type dayClosureHandler struct {
	service ports.AuthorizedDayClosureService
}

var _ DayClosureHandler = (*dayClosureHandler)(nil)

func NewDayClosureHandler(service ports.AuthorizedDayClosureService) *dayClosureHandler {
	return &dayClosureHandler{
		service: service,
	}
}

type closeDayRequest struct {
	Date string `json:"date"`
}

func (h *dayClosureHandler) Close(r *http.Request) (*dayClosureResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[closeDayRequest](r)
	if err != nil {
		return nil, 0, err
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	err = h.service.Close(r.Context(), userId, date)
	if err != nil {
		return nil, 0, err
	}
	response := &dayClosureResponse{
		Date:   request.Date,
		Closed: true,
	}
	return response, http.StatusOK, nil
}

func (h *dayClosureHandler) Open(r *http.Request) (*api.NoResponse, int, error) {
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
	err = h.service.Open(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}

func (h *dayClosureHandler) GetStatus(r *http.Request) (*dayClosureResponse, int, error) {
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
	closed, err := h.service.IsClosed(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	response := &dayClosureResponse{
		Date:   date.Format("2006-01-02"),
		Closed: closed,
	}
	return response, http.StatusOK, nil
}
