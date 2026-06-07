package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/weight/ports"
)

type WeightEntryHandler interface {
	Create(r *http.Request) (*weightEntryResponse, int, error)
	GetById(r *http.Request) (*weightEntryResponse, int, error)
	List(r *http.Request) (*weightEntryPage, int, error)
	Update(r *http.Request) (*weightEntryResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
}

type weightEntryHandler struct {
	service ports.AuthorizedWeightEntryService
}

var _ WeightEntryHandler = (*weightEntryHandler)(nil)

func NewWeightEntryHandler(service ports.AuthorizedWeightEntryService) *weightEntryHandler {
	return &weightEntryHandler{
		service: service,
	}
}

func (h *weightEntryHandler) Create(r *http.Request) (*weightEntryResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[createWeightEntryRequest](r)
	if err != nil {
		return nil, 0, err
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	params := ports.CreateParams{
		Date:              date,
		WeightKg:          request.WeightKg,
		BodyFatPercentage: request.BodyFatPercentage,
		Notes:             request.Notes,
	}
	entry, err := h.service.Create(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return weightEntryFromDomain(entry), http.StatusCreated, nil
}

func (h *weightEntryHandler) GetById(r *http.Request) (*weightEntryResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	entry, err := h.service.GetById(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return weightEntryFromDomain(entry), http.StatusOK, nil
}

func (h *weightEntryHandler) List(r *http.Request) (*weightEntryPage, int, error) {
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
	params := ports.ListParams{
		PaginationParams: api.PaginationFromRequest(r),
		From:             from,
		To:               to,
	}
	page, err := h.service.List(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return newWeightEntryPage(page), http.StatusOK, nil
}

func (h *weightEntryHandler) Update(r *http.Request) (*weightEntryResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateWeightEntryRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpdateParams{
		WeightKg:          request.WeightKg,
		BodyFatPercentage: request.BodyFatPercentage,
		Notes:             request.Notes,
	}
	if request.Date != nil {
		date, err := time.Parse("2006-01-02", *request.Date)
		if err != nil {
			return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
		}
		params.Date = &date
	}
	entry, err := h.service.Update(r.Context(), userId, id, params)
	if err != nil {
		return nil, 0, err
	}
	return weightEntryFromDomain(entry), http.StatusOK, nil
}

func (h *weightEntryHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	err = h.service.Delete(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}
