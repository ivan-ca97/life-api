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
	ListCommunity(r *http.Request) (*foodPage, int, error)
	Copy(r *http.Request) (*foodResponse, int, error)
	Impact(r *http.Request) (*impactResponse, int, error)
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
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
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
	portions := make([]ports.PortionParam, len(request.Portions))
	for i, p := range request.Portions {
		portions[i] = ports.PortionParam{
			Name:           p.Name,
			BaseEquivalent: p.BaseEquivalent,
		}
	}
	params := ports.CreateParams{
		Name:                request.Name,
		PhotoUrl:            request.PhotoUrl,
		DefaultCalories:     request.DefaultCalories,
		DefaultProteinGrams: request.DefaultProteinGrams,
		DefaultCarbsGrams:   request.DefaultCarbsGrams,
		DefaultFatGrams:     request.DefaultFatGrams,
		DefaultFiberGrams:   request.DefaultFiberGrams,
		MeasurementType:     request.MeasurementType,
		BaseQuantity:        baseQuantity,
		BaseUnit:            request.BaseUnit,
		Public:              request.Public,
		Tags:                request.Tags,
		Ingredients:         request.Ingredients,
		Conversions:         conversionsParamFromRequest(request.Conversions),
		Portions:            portions,
	}
	food, err := h.service.Create(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusCreated, nil
}

func (h *foodHandler) GetById(r *http.Request) (*foodResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	food, err := h.service.GetById(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusOK, nil
}

func (h *foodHandler) List(r *http.Request) (*foodPage, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	var query *string
	if q := r.URL.Query().Get("q"); q != "" {
		query = &q
	}
	var tag *string
	if t := r.URL.Query().Get("tag"); t != "" {
		tag = &t
	}
	var sort *string
	if s := r.URL.Query().Get("sort"); s != "" {
		sort = &s
	}
	params := ports.ListParams{
		PaginationParams: api.PaginationFromRequest(r),
		Query:            query,
		Tag:              tag,
		Sort:             sort,
	}
	page, err := h.service.List(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return newFoodPage(page), http.StatusOK, nil
}

func (h *foodHandler) Update(r *http.Request) (*foodResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateFoodRequest](r)
	if err != nil {
		return nil, 0, err
	}
	var portions *[]ports.PortionParam
	if request.Portions != nil {
		pp := make([]ports.PortionParam, len(*request.Portions))
		for i, p := range *request.Portions {
			pp[i] = ports.PortionParam{Name: p.Name, BaseEquivalent: p.BaseEquivalent}
		}
		portions = &pp
	}
	params := ports.UpdateParams{
		Name:                request.Name,
		PhotoUrl:            request.PhotoUrl,
		DefaultCalories:     request.DefaultCalories,
		DefaultProteinGrams: request.DefaultProteinGrams,
		DefaultCarbsGrams:   request.DefaultCarbsGrams,
		DefaultFatGrams:     request.DefaultFatGrams,
		DefaultFiberGrams:   request.DefaultFiberGrams,
		MeasurementType:     request.MeasurementType,
		BaseQuantity:        request.BaseQuantity,
		BaseUnit:            request.BaseUnit,
		Public:              request.Public,
		Tags:                request.Tags,
		Ingredients:         request.Ingredients,
		Conversions:         conversionsParamFromRequest(request.Conversions),
		Portions:            portions,
	}
	food, err := h.service.Update(r.Context(), userId, id, params)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusOK, nil
}

func (h *foodHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
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

func (h *foodHandler) Frequency(r *http.Request) (*frequencyResponse, int, error) {
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
	var tag *string
	if t := r.URL.Query().Get("tag"); t != "" {
		tag = &t
	}
	params := ports.FrequencyParams{
		From: from,
		To:   to,
		Tag:  tag,
	}
	results, err := h.service.Frequency(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return newFrequencyResponse(results), http.StatusOK, nil
}

func (h *foodHandler) ListIngredients(r *http.Request) (*ingredientsListResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	var query *string
	if q := r.URL.Query().Get("q"); q != "" {
		query = &q
	}
	ingredients, err := h.service.ListIngredients(r.Context(), userId, query)
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
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	food, err := h.service.GetById(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	metric := domain.MetricUnitsForDimension(domain.UnitDimension(food.MeasurementType))
	var extraUnits []string
	if food.VolumeConversion != nil {
		if domain.UnitDimension(food.MeasurementType) == domain.DimensionMass {
			extraUnits = append(extraUnits, domain.MetricUnitsForDimension(domain.DimensionVolume)...)
		} else {
			extraUnits = append(extraUnits, domain.MetricUnitsForDimension(domain.DimensionMass)...)
		}
	}
	if food.UnitConversion != nil {
		extraUnits = append(extraUnits, "u")
	}
	for _, p := range food.Portions {
		extraUnits = append(extraUnits, p.Name)
	}
	return &foodUnitsResponse{
		Metric:      metric,
		Conversions: extraUnits,
	}, http.StatusOK, nil
}

func conversionsParamFromRequest(r *conversionsRequest) *ports.ConversionsParam {
	if r == nil {
		return nil
	}
	p := &ports.ConversionsParam{}
	if r.VolumeConversion != nil {
		p.VolumeConversion = &ports.VolumeConversionParam{
			GramsPerMl: r.VolumeConversion.GramsPerMl,
			Note:       r.VolumeConversion.Note,
		}
	}
	if r.UnitConversion != nil {
		p.UnitConversion = &ports.UnitConversionParam{
			BaseEquivalent: r.UnitConversion.BaseEquivalent,
			Note:           r.UnitConversion.Note,
		}
	}
	return p
}

func (h *foodHandler) IngredientFrequency(r *http.Request) (*ingredientFrequencyResponse, int, error) {
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
	params := ports.IngredientFrequencyParams{
		From: from,
		To:   to,
	}
	results, err := h.service.IngredientFrequency(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return newIngredientFrequencyResponse(results), http.StatusOK, nil
}

func (h *foodHandler) ListCommunity(r *http.Request) (*foodPage, int, error) {
	var query *string
	q := r.URL.Query().Get("q")
	if q != "" {
		query = &q
	}
	params := ports.CommunityListParams{
		PaginationParams: api.PaginationFromRequest(r),
		Query:            query,
	}
	page, err := h.service.ListCommunity(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return newFoodPage(page), http.StatusOK, nil
}

func (h *foodHandler) Copy(r *http.Request) (*foodResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	food, err := h.service.Copy(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return foodFromDomain(food), http.StatusCreated, nil
}

func (h *foodHandler) Impact(r *http.Request) (*impactResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	result, err := h.service.Impact(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return newImpactResponse(result), http.StatusOK, nil
}
