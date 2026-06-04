package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/food/domain"
	"github.com/ivan-ca97/life/internal/features/food/ports"
)

type FoodHandler interface {
	Create(r *http.Request) (*foodResponse, int, error)
	GetById(r *http.Request) (*foodResponse, int, error)
	List(r *http.Request) (*foodPage, int, error)
	Update(r *http.Request) (*foodResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
	Frequency(r *http.Request) (*frequencyResponse, int, error)
	IngredientFrequency(r *http.Request) (*ingredientFrequencyResponse, int, error)
	ListIngredients(r *http.Request) (*ingredientsListResponse, int, error)
	ListUnits(r *http.Request) (*unitsListResponse, int, error)
	ListFoodUnits(r *http.Request) (*foodUnitsResponse, int, error)
}

type foodHandler struct {
	service ports.AuthorizedFoodService
}

var _ FoodHandler = (*foodHandler)(nil)

func NewFoodHandler(service ports.AuthorizedFoodService) *foodHandler {
	return &foodHandler{
		service: service,
	}
}

func (h *foodHandler) Create(r *http.Request) (*foodResponse, int, error) {
	request, err := api.DecodeBody[createFoodRequest](r)
	if err != nil {
		return nil, 0, err
	}
	if request.Name == "" {
		return nil, 0, cerr.NewBadRequestError("name is required")
	}
	baseQuantity := 1.0
	if request.BaseQuantity != nil {
		baseQuantity = *request.BaseQuantity
	}
	conversions := make([]ports.ConversionParam, len(request.Conversions))
	for i, c := range request.Conversions {
		conversions[i] = ports.ConversionParam{
			Unit:           c.Unit,
			BaseEquivalent: c.BaseEquivalent,
			Inverse:        c.Inverse,
			Note:           c.Note,
		}
	}
	params := ports.CreateParams{
		Name:                request.Name,
		DefaultCalories:     request.DefaultCalories,
		DefaultProteinGrams: request.DefaultProteinGrams,
		DefaultCarbsGrams:   request.DefaultCarbsGrams,
		DefaultFatGrams:     request.DefaultFatGrams,
		DefaultFiberGrams:   request.DefaultFiberGrams,
		MeasurementType:     request.MeasurementType,
		BaseQuantity:        baseQuantity,
		BaseUnit:            request.BaseUnit,
		Tags:                request.Tags,
		Ingredients:         request.Ingredients,
		Conversions:         conversions,
	}
	food, err := h.service.Create(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusCreated, nil
}

func (h *foodHandler) GetById(r *http.Request) (*foodResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	food, err := h.service.GetById(r.Context(), id)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusOK, nil
}

func (h *foodHandler) List(r *http.Request) (*foodPage, int, error) {
	var query *string
	if q := r.URL.Query().Get("q"); q != "" {
		query = &q
	}
	var tag *string
	if t := r.URL.Query().Get("tag"); t != "" {
		tag = &t
	}
	params := ports.ListParams{
		PaginationParams: api.PaginationFromRequest(r),
		Query:            query,
		Tag:              tag,
	}
	page, err := h.service.List(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return newFoodPage(page), http.StatusOK, nil
}

func (h *foodHandler) Update(r *http.Request) (*foodResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateFoodRequest](r)
	if err != nil {
		return nil, 0, err
	}
	var conversions *[]ports.ConversionParam
	if request.Conversions != nil {
		convs := make([]ports.ConversionParam, len(*request.Conversions))
		for i, c := range *request.Conversions {
			convs[i] = ports.ConversionParam{
				Unit:           c.Unit,
				BaseEquivalent: c.BaseEquivalent,
				Inverse:        c.Inverse,
				Note:           c.Note,
			}
		}
		conversions = &convs
	}
	params := ports.UpdateParams{
		Name:                request.Name,
		DefaultCalories:     request.DefaultCalories,
		DefaultProteinGrams: request.DefaultProteinGrams,
		DefaultCarbsGrams:   request.DefaultCarbsGrams,
		DefaultFatGrams:     request.DefaultFatGrams,
		DefaultFiberGrams:   request.DefaultFiberGrams,
		MeasurementType:     request.MeasurementType,
		BaseQuantity:        request.BaseQuantity,
		BaseUnit:            request.BaseUnit,
		Tags:                request.Tags,
		Ingredients:         request.Ingredients,
		Conversions:         conversions,
	}
	food, err := h.service.Update(r.Context(), id, params)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusOK, nil
}

func (h *foodHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	err = h.service.Delete(r.Context(), id)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}

func (h *foodHandler) Frequency(r *http.Request) (*frequencyResponse, int, error) {
	from, err := api.QueryParamDate(r, "from")
	if err != nil {
		return nil, 0, err
	}
	to, err := api.QueryParamDate(r, "to")
	if err != nil {
		return nil, 0, err
	}
	var tag *string
	if t := r.URL.Query().Get("tag"); t != "" {
		tag = &t
	}
	params := ports.FrequencyParams{
		From: from,
		To:   to,
		Tag:  tag,
	}
	results, err := h.service.Frequency(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return newFrequencyResponse(results), http.StatusOK, nil
}

func (h *foodHandler) ListIngredients(r *http.Request) (*ingredientsListResponse, int, error) {
	var query *string
	if q := r.URL.Query().Get("q"); q != "" {
		query = &q
	}
	ingredients, err := h.service.ListIngredients(r.Context(), query)
	if err != nil {
		return nil, 0, err
	}
	items := make([]ingredientResponse, len(ingredients))
	for i, ing := range ingredients {
		items[i] = ingredientResponse{Id: ing.Id, Name: ing.Name}
	}
	return &ingredientsListResponse{Items: items}, http.StatusOK, nil
}

func (h *foodHandler) ListUnits(_ *http.Request) (*unitsListResponse, int, error) {
	return &unitsListResponse{
		Mass:   domain.MetricUnitsForDimension(domain.DimensionMass),
		Volume: domain.MetricUnitsForDimension(domain.DimensionVolume),
		Unit:   domain.MetricUnitsForDimension(domain.DimensionUnit),
	}, http.StatusOK, nil
}

func (h *foodHandler) ListFoodUnits(r *http.Request) (*foodUnitsResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	food, err := h.service.GetById(r.Context(), id)
	if err != nil {
		return nil, 0, err
	}
	metric := domain.MetricUnitsForDimension(domain.UnitDimension(food.MeasurementType))
	conversions := make([]string, len(food.Conversions))
	for i, c := range food.Conversions {
		conversions[i] = c.Unit
	}
	return &foodUnitsResponse{
		Metric:      metric,
		Conversions: conversions,
	}, http.StatusOK, nil
}

func (h *foodHandler) IngredientFrequency(r *http.Request) (*ingredientFrequencyResponse, int, error) {
	from, err := api.QueryParamDate(r, "from")
	if err != nil {
		return nil, 0, err
	}
	to, err := api.QueryParamDate(r, "to")
	if err != nil {
		return nil, 0, err
	}
	params := ports.IngredientFrequencyParams{
		From: from,
		To:   to,
	}
	results, err := h.service.IngredientFrequency(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return newIngredientFrequencyResponse(results), http.StatusOK, nil
}
