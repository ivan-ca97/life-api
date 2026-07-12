package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/validate"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

// --- request ---

type upsertCorrectionRequest struct {
	Date            string   `json:"date"`
	Calories        *float64 `json:"calories"`
	ProteinGrams    *float64 `json:"protein_grams"`
	CarbsGrams      *float64 `json:"carbs_grams"`
	FatGrams        *float64 `json:"fat_grams"`
	FiberGrams      *float64 `json:"fiber_grams"`
	CaloriesBurned  *float64 `json:"calories_burned"`
	Steps           *int     `json:"steps"`
	DurationSeconds *int     `json:"duration_seconds"`
	DistanceMeters  *float64 `json:"distance_meters"`
	Notes           string   `json:"notes"`
}

// --- response ---

type correctionResponse struct {
	Date            string   `json:"date"`
	Calories        *float64 `json:"calories,omitempty"`
	ProteinGrams    *float64 `json:"protein_grams,omitempty"`
	CarbsGrams      *float64 `json:"carbs_grams,omitempty"`
	FatGrams        *float64 `json:"fat_grams,omitempty"`
	FiberGrams      *float64 `json:"fiber_grams,omitempty"`
	CaloriesBurned  *float64 `json:"calories_burned,omitempty"`
	Steps           *int     `json:"steps,omitempty"`
	DurationSeconds *int     `json:"duration_seconds,omitempty"`
	DistanceMeters  *float64 `json:"distance_meters,omitempty"`
	Notes           string   `json:"notes"`
}

func correctionFromDomain(c *domain.Correction) *correctionResponse {
	return &correctionResponse{
		Date:            c.Date.Format("2006-01-02"),
		Calories:        c.Calories,
		ProteinGrams:    c.ProteinGrams,
		CarbsGrams:      c.CarbsGrams,
		FatGrams:        c.FatGrams,
		FiberGrams:      c.FiberGrams,
		CaloriesBurned:  c.CaloriesBurned,
		Steps:           c.Steps,
		DurationSeconds: c.DurationSeconds,
		DistanceMeters:  c.DistanceMeters,
		Notes:           c.Notes,
	}
}

// --- handler ---

type CorrectionHandler interface {
	GetCorrection(r *http.Request) (*correctionResponse, int, error)
	UpsertCorrection(r *http.Request) (*correctionResponse, int, error)
	DeleteCorrection(r *http.Request) (*api.NoResponse, int, error)
}

type correctionHandler struct {
	service ports.AuthorizedCorrectionService
}

var _ CorrectionHandler = (*correctionHandler)(nil)

func NewCorrectionHandler(service ports.AuthorizedCorrectionService) *correctionHandler {
	return &correctionHandler{service: service}
}

func (h *correctionHandler) GetCorrection(r *http.Request) (*correctionResponse, int, error) {
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
	correction, err := h.service.GetCorrection(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	if correction == nil {
		correction = &domain.Correction{
			Date: *date,
		}
	}
	return correctionFromDomain(correction), http.StatusOK, nil
}

func (h *correctionHandler) UpsertCorrection(r *http.Request) (*correctionResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[upsertCorrectionRequest](r)
	if err != nil {
		return nil, 0, err
	}
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		return nil, 0, cerr.NewBadRequestError("invalid date format, expected YYYY-MM-DD")
	}
	err = validate.NonNegativePtr(request.Calories, "calories")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativePtr(request.ProteinGrams, "protein_grams")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativePtr(request.CarbsGrams, "carbs_grams")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativePtr(request.FatGrams, "fat_grams")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativePtr(request.FiberGrams, "fiber_grams")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativePtr(request.CaloriesBurned, "calories_burned")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativeIntPtr(request.Steps, "steps")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativeIntPtr(request.DurationSeconds, "duration_seconds")
	if err != nil {
		return nil, 0, err
	}
	err = validate.NonNegativePtr(request.DistanceMeters, "distance_meters")
	if err != nil {
		return nil, 0, err
	}
	correction := &domain.Correction{
		Date:            date,
		Calories:        request.Calories,
		ProteinGrams:    request.ProteinGrams,
		CarbsGrams:      request.CarbsGrams,
		FatGrams:        request.FatGrams,
		FiberGrams:      request.FiberGrams,
		CaloriesBurned:  request.CaloriesBurned,
		Steps:           request.Steps,
		DurationSeconds: request.DurationSeconds,
		DistanceMeters:  request.DistanceMeters,
		Notes:           request.Notes,
	}
	err = h.service.UpsertCorrection(r.Context(), userId, correction)
	if err != nil {
		return nil, 0, err
	}
	updated, err := h.service.GetCorrection(r.Context(), userId, date)
	if err != nil {
		return nil, 0, err
	}
	if updated == nil {
		updated = &domain.Correction{
			Date: date,
		}
	}
	return correctionFromDomain(updated), http.StatusOK, nil
}

func (h *correctionHandler) DeleteCorrection(r *http.Request) (*api.NoResponse, int, error) {
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
	err = h.service.DeleteCorrection(r.Context(), userId, *date)
	if err != nil {
		return nil, 0, err
	}
	return nil, http.StatusNoContent, nil
}
