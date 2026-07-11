package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/measurements/ports"
)

type MeasurementHandler interface {
	Upsert(r *http.Request) (*measurementResponse, int, error)
	GetByDate(r *http.Request) (*measurementResponse, int, error)
	List(r *http.Request) (*measurementListResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
}

type measurementHandler struct {
	service ports.AuthorizedMeasurementService
}

var _ MeasurementHandler = (*measurementHandler)(nil)

func NewMeasurementHandler(service ports.AuthorizedMeasurementService) *measurementHandler {
	return &measurementHandler{service: service}
}

func (h *measurementHandler) Upsert(r *http.Request) (*measurementResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, measureType, err := parseDateAndType(r)
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[upsertMeasurementRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpsertParams{
		Value: request.Value,
		Notes: request.Notes,
	}
	m, err := h.service.Upsert(r.Context(), userId, date, measureType, params)
	if err != nil {
		return nil, 0, err
	}
	return measurementFromDomain(m), http.StatusOK, nil
}

func (h *measurementHandler) GetByDate(r *http.Request) (*measurementResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, measureType, err := parseDateAndType(r)
	if err != nil {
		return nil, 0, err
	}
	m, err := h.service.GetByDate(r.Context(), userId, date, measureType)
	if err != nil {
		return nil, 0, err
	}
	return measurementFromDomain(m), http.StatusOK, nil
}

func (h *measurementHandler) List(r *http.Request) (*measurementListResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	from, err := api.QueryParamDate(r, "from")
	if err != nil {
		return nil, 0, err
	}
	to, err := api.QueryParamDate(r, "to")
	if err != nil {
		return nil, 0, err
	}
	var measureType *string
	if t := r.URL.Query().Get("type"); t != "" {
		measureType = &t
	}
	entries, err := h.service.List(r.Context(), userId, ports.ListParams{
		From: from,
		To:   to,
		Type: measureType,
	})
	if err != nil {
		return nil, 0, err
	}
	items := make([]measurementResponse, len(entries))
	for i, e := range entries {
		items[i] = *measurementFromDomain(&e)
	}
	return &measurementListResponse{Items: items}, http.StatusOK, nil
}

func (h *measurementHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, measureType, err := parseDateAndType(r)
	if err != nil {
		return nil, 0, err
	}
	if err := h.service.Delete(r.Context(), userId, date, measureType); err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}

func parseDateAndType(r *http.Request) (time.Time, string, error) {
	raw := api.PathParam(r, "date")
	date, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, "", cerr.NewBadRequestError("date must be in YYYY-MM-DD format")
	}
	measureType := api.PathParam(r, "type")
	if measureType == "" {
		return time.Time{}, "", cerr.NewBadRequestError("type is required")
	}
	return date, measureType, nil
}
