package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/steps/ports"
)

type StepsHandler interface {
	Upsert(r *http.Request) (*stepsResponse, int, error)
	GetByDate(r *http.Request) (*stepsResponse, int, error)
	List(r *http.Request) (*stepsListResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
}

type stepsHandler struct {
	service      ports.AuthorizedStepsService
	weightLookup ports.WeightLookup
}

var _ StepsHandler = (*stepsHandler)(nil)

func NewStepsHandler(service ports.AuthorizedStepsService, weightLookup ports.WeightLookup) *stepsHandler {
	return &stepsHandler{service: service, weightLookup: weightLookup}
}

func (h *stepsHandler) Upsert(r *http.Request) (*stepsResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := parseDate(r)
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[upsertStepsRequest](r)
	if err != nil {
		return nil, 0, err
	}
	entry, err := h.service.Upsert(r.Context(), userId, date, ports.UpsertParams{
		Steps:  request.Steps,
		Source: request.Source,
	})
	if err != nil {
		return nil, 0, err
	}
	return stepsFromDomain(entry, h.lookupWeight(userId)), http.StatusOK, nil
}

func (h *stepsHandler) GetByDate(r *http.Request) (*stepsResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := parseDate(r)
	if err != nil {
		return nil, 0, err
	}
	entry, err := h.service.GetByDate(r.Context(), userId, date)
	if err != nil {
		return nil, 0, err
	}
	return stepsFromDomain(entry, h.lookupWeight(userId)), http.StatusOK, nil
}

func (h *stepsHandler) List(r *http.Request) (*stepsListResponse, int, error) {
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
	entries, err := h.service.List(r.Context(), userId, ports.ListParams{From: from, To: to})
	if err != nil {
		return nil, 0, err
	}
	weightKg := h.lookupWeight(userId)
	items := make([]stepsResponse, len(entries))
	for i, e := range entries {
		items[i] = *stepsFromDomain(&e, weightKg)
	}
	return &stepsListResponse{Items: items}, http.StatusOK, nil
}

func (h *stepsHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := parseDate(r)
	if err != nil {
		return nil, 0, err
	}
	if err := h.service.Delete(r.Context(), userId, date); err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}

func (h *stepsHandler) lookupWeight(userId uuid.UUID) *float64 {
	kg, found, err := h.weightLookup.LatestWeightKg(userId)
	if err != nil || !found {
		return nil
	}
	return &kg
}

func parseDate(r *http.Request) (time.Time, error) {
	raw := api.PathParam(r, "date")
	date, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, cerr.NewBadRequestError("date must be in YYYY-MM-DD format")
	}
	return date, nil
}
