package handler

import (
	"net/http"
	"time"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/goal/ports"
)

type GoalHandler interface {
	GetCurrent(r *http.Request) (*goalResponse, int, error)
	Upsert(r *http.Request) (*goalResponse, int, error)
	GetProgress(r *http.Request) (*goalProgressResponse, int, error)
}

type goalHandler struct {
	service ports.AuthorizedGoalService
}

var _ GoalHandler = (*goalHandler)(nil)

func NewGoalHandler(service ports.AuthorizedGoalService) *goalHandler {
	return &goalHandler{
		service: service,
	}
}

func (h *goalHandler) GetCurrent(r *http.Request) (*goalResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	goal, err := h.service.GetCurrent(r.Context(), userId)
	if err != nil {
		return nil, 0, err
	}
	return goalFromDomain(goal), http.StatusOK, nil
}

func (h *goalHandler) Upsert(r *http.Request) (*goalResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[upsertGoalRequest](r)
	if err != nil {
		return nil, 0, err
	}
	params := ports.UpsertParams{
		DailyCalories:        request.DailyCalories,
		DailyProteinGrams:    request.DailyProteinGrams,
		DailyCarbsGrams:      request.DailyCarbsGrams,
		DailyFatGrams:        request.DailyFatGrams,
		DailyFiberGrams:      request.DailyFiberGrams,
		DailySteps:           request.DailySteps,
		DailyExerciseMinutes: request.DailyExerciseMinutes,
		TargetWeightKg:       request.TargetWeightKg,
	}
	if request.StartedAt != nil {
		startedAt, err := time.Parse(time.RFC3339, *request.StartedAt)
		if err != nil {
			return nil, 0, cerr.NewBadRequestError("invalid started_at format, expected RFC3339")
		}
		params.StartedAt = &startedAt
	}
	goal, err := h.service.Upsert(r.Context(), userId, params)
	if err != nil {
		return nil, 0, err
	}
	return goalFromDomain(goal), http.StatusOK, nil
}

func (h *goalHandler) GetProgress(r *http.Request) (*goalProgressResponse, int, error) {
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
	if to.Before(*from) {
		return nil, 0, cerr.NewBadRequestError("to must not be before from")
	}
	progress, err := h.service.GetProgress(r.Context(), userId, *from, *to)
	if err != nil {
		return nil, 0, err
	}
	return goalProgressFromDomain(progress), http.StatusOK, nil
}
