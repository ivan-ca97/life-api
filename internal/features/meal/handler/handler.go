package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
	"github.com/ivan-ca97/life/internal/features/meal/ports"
)

type MealHandler interface {
	Create(r *http.Request) (*mealResponse, int, error)
	GetById(r *http.Request) (*mealResponse, int, error)
	List(r *http.Request) (*mealPage, int, error)
	Update(r *http.Request) (*mealResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
	ListTypes(r *http.Request) (*mealTypesResponse, int, error)
	PreviewNutrition(r *http.Request) (*nutritionPreviewResponse, int, error)
}

type mealHandler struct {
	service ports.AuthorizedMealService
}

var _ MealHandler = (*mealHandler)(nil)

func NewMealHandler(service ports.AuthorizedMealService) *mealHandler {
	return &mealHandler{
		service: service,
	}
}

func (h *mealHandler) Create(r *http.Request) (*mealResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[createMealRequest](r)
	if err != nil {
		return nil, 0, err
	}
	if !request.hasContent() {
		return nil, 0, cerr.NewBadRequestError("meal must have at least one detail: name, notes, photo, items, or nutrition values")
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	items := make([]ports.ItemParam, len(request.Items))
	for i, item := range request.Items {
		method := domain.MeasurementMethod(item.MeasurementMethod)
		if !domain.IsValidMeasurementMethod(method) {
			return nil, 0, cerr.NewBadRequestError("invalid measurement_method: " + item.MeasurementMethod)
		}
		items[i] = ports.ItemParam{
			FoodId:            item.FoodId,
			Quantity:          item.Quantity,
			Unit:              item.Unit,
			Notes:             item.Notes,
			MeasurementMethod: method,
		}
	}
	params := ports.CreateParams{
		Date:         date,
		Type:         request.Type,
		Name:         request.Name,
		PhotoUrl:     request.PhotoUrl,
		EatenAt:      request.EatenAt,
		Calories:     request.Calories,
		ProteinGrams: request.ProteinGrams,
		CarbsGrams:   request.CarbsGrams,
		FatGrams:     request.FatGrams,
		FiberGrams:   request.FiberGrams,
		Tags:         request.Tags,
		Items:        items,
		Notes:        request.Notes,
	}
	meal, err := h.service.Create(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return mealFromDomain(meal), http.StatusCreated, nil
}

func (h *mealHandler) GetById(r *http.Request) (*mealResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	meal, err := h.service.GetById(r.Context(), userId, id)
	if err != nil {
		return nil, 0, err
	}
	return mealFromDomain(meal), http.StatusOK, nil
}

func (h *mealHandler) List(r *http.Request) (*mealPage, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	date, err := api.QueryParamDate(r, "date")
	if err != nil {
		return nil, 0, err
	}
	params := ports.ListParams{
		PaginationParams: api.PaginationFromRequest(r),
		Date:             date,
	}
	page, err := h.service.List(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return newMealPage(page), http.StatusOK, nil
}

func (h *mealHandler) Update(r *http.Request) (*mealResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateMealRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpdateParams{
		Type:         request.Type,
		Name:         request.Name,
		PhotoUrl:     request.PhotoUrl,
		EatenAt:      request.EatenAt,
		Calories:     request.Calories,
		ProteinGrams: request.ProteinGrams,
		CarbsGrams:   request.CarbsGrams,
		FatGrams:     request.FatGrams,
		FiberGrams:   request.FiberGrams,
		Tags:         request.Tags,
		Notes:        request.Notes,
	}
	if request.Items != nil {
		items := make([]ports.ItemParam, len(*request.Items))
		for i, item := range *request.Items {
			method := domain.MeasurementMethod(item.MeasurementMethod)
			if !domain.IsValidMeasurementMethod(method) {
				return nil, 0, cerr.NewBadRequestError("invalid measurement_method: " + item.MeasurementMethod)
			}
			items[i] = ports.ItemParam{
				FoodId:            item.FoodId,
				Quantity:          item.Quantity,
				Unit:              item.Unit,
				Notes:             item.Notes,
				MeasurementMethod: method,
			}
		}
		params.Items = &items
	}
	if request.Date != nil {
		date, err := time.Parse("2006-01-02", *request.Date)
		if err != nil {
			return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
		}
		params.Date = &date
	}
	meal, err := h.service.Update(r.Context(), userId, id, params)
	if err != nil {
		return nil, 0, err
	}
	return mealFromDomain(meal), http.StatusOK, nil
}

func (h *mealHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
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

func (h *mealHandler) PreviewNutrition(r *http.Request) (*nutritionPreviewResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[previewNutritionRequest](r)
	if err != nil {
		return nil, 0, err
	}
	items := make([]ports.ItemParam, len(request.Items))
	for i, item := range request.Items {
		items[i] = ports.ItemParam{
			FoodId:   item.FoodId,
			Quantity: item.Quantity,
			Unit:     item.Unit,
		}
	}
	preview, err := h.service.PreviewNutrition(r.Context(), userId, items)
	if err != nil {
		return nil, 0, err
	}
	respItems := make([]nutritionPreviewItemResponse, len(preview.Items))
	for i, pi := range preview.Items {
		respItems[i] = nutritionPreviewItemResponse{
			FoodId:       pi.FoodId,
			Calories:     pi.Calories,
			ProteinGrams: pi.ProteinGrams,
			CarbsGrams:   pi.CarbsGrams,
			FatGrams:     pi.FatGrams,
			FiberGrams:   pi.FiberGrams,
		}
	}
	response := &nutritionPreviewResponse{
		Calories:     preview.Calories,
		ProteinGrams: preview.ProteinGrams,
		CarbsGrams:   preview.CarbsGrams,
		FatGrams:     preview.FatGrams,
		FiberGrams:   preview.FiberGrams,
		Items:        respItems,
	}
	return response, http.StatusOK, nil
}

func (h *mealHandler) ListTypes(r *http.Request) (*mealTypesResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	hour, err := api.QueryParamInt(r, "hour")
	if err != nil {
		return nil, 0, err
	}
	if hour != nil && (*hour < 0 || *hour > 23) {
		return nil, 0, cerr.NewBadRequestError("hour must be between 0 and 23")
	}
	mealTypes, err := h.service.ListTypes(r.Context(), userId, hour)
	if err != nil {
		return nil, 0, err
	}
	response := &mealTypesResponse{
		Types: mealTypes,
	}
	return response, http.StatusOK, nil
}
