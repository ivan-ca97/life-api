package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/exercise/ports"
)

type ExerciseHandler interface {
	Create(r *http.Request) (*exerciseResponse, int, error)
	GetById(r *http.Request) (*exerciseResponse, int, error)
	List(r *http.Request) (*exercisePage, int, error)
	Update(r *http.Request) (*exerciseResponse, int, error)
	Delete(r *http.Request) (*api.NoResponse, int, error)
}

type exerciseHandler struct {
	service ports.AuthorizedExerciseService
}

var _ ExerciseHandler = (*exerciseHandler)(nil)

func NewExerciseHandler(service ports.AuthorizedExerciseService) *exerciseHandler {
	return &exerciseHandler{
		service: service,
	}
}

func (h *exerciseHandler) Create(r *http.Request) (*exerciseResponse, int, error) {
	request, err := api.DecodeBody[createExerciseRequest](r)
	if err != nil {
		return nil, 0, err
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	params := ports.CreateParams{
		Date:                    date,
		Type:                    request.Type,
		Name:                    request.Name,
		StartedAt:               request.StartedAt,
		DurationSeconds:         request.DurationSeconds,
		EstimatedCaloriesBurned: request.EstimatedCaloriesBurned,
		Steps:                   request.Steps,
		DistanceMeters:          request.DistanceMeters,
		MaxSpeedKmh:             request.MaxSpeedKmh,
		ElevationGainMeters:     request.ElevationGainMeters,
		AverageHeartRate:        request.AverageHeartRate,
		MaxHeartRate:            request.MaxHeartRate,
		TotalVolumeKg:           request.TotalVolumeKg,
		TotalSets:               request.TotalSets,
		Tags:                    request.Tags,
		Notes:                   request.Notes,
	}
	exercise, err := h.service.Create(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return exerciseFromDomain(exercise), http.StatusCreated, nil
}

func (h *exerciseHandler) GetById(r *http.Request) (*exerciseResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	exercise, err := h.service.GetById(r.Context(), id)
	if err != nil {
		return nil, 0, err
	}
	return exerciseFromDomain(exercise), http.StatusOK, nil
}

func (h *exerciseHandler) List(r *http.Request) (*exercisePage, int, error) {
	date, err := api.QueryParamDate(r, "date")
	if err != nil {
		return nil, 0, err
	}
	params := ports.ListParams{
		PaginationParams: api.PaginationFromRequest(r),
		Date:             date,
	}
	page, err := h.service.List(r.Context(), params)
	if err != nil {
		return nil, 0, err
	}
	return newExercisePage(page), http.StatusOK, nil
}

func (h *exerciseHandler) Update(r *http.Request) (*exerciseResponse, int, error) {
	id, err := api.PathParamUUID(r, "id")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[updateExerciseRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpdateParams{
		Type:                    request.Type,
		Name:                    request.Name,
		StartedAt:               request.StartedAt,
		DurationSeconds:         request.DurationSeconds,
		EstimatedCaloriesBurned: request.EstimatedCaloriesBurned,
		Steps:                   request.Steps,
		DistanceMeters:          request.DistanceMeters,
		MaxSpeedKmh:             request.MaxSpeedKmh,
		ElevationGainMeters:     request.ElevationGainMeters,
		AverageHeartRate:        request.AverageHeartRate,
		MaxHeartRate:            request.MaxHeartRate,
		TotalVolumeKg:           request.TotalVolumeKg,
		TotalSets:               request.TotalSets,
		Tags:                    request.Tags,
		Notes:                   request.Notes,
	}
	if request.Date != nil {
		date, err := time.Parse("2006-01-02", *request.Date)
		if err != nil {
			return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
		}
		params.Date = &date
	}
	exercise, err := h.service.Update(r.Context(), id, params)
	if err != nil {
		return nil, 0, err
	}
	return exerciseFromDomain(exercise), http.StatusOK, nil
}

func (h *exerciseHandler) Delete(r *http.Request) (*api.NoResponse, int, error) {
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
