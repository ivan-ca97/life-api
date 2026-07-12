package handler

import (
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

type MealAIHandler interface {
	Estimate(r *http.Request) (*estimateResponse, int, error)
}

type mealAIHandler struct {
	useCase ports.MealEstimationUseCase
}

var _ MealAIHandler = (*mealAIHandler)(nil)

func NewMealAIHandler(useCase ports.MealEstimationUseCase) *mealAIHandler {
	return &mealAIHandler{useCase: useCase}
}

func (h *mealAIHandler) Estimate(r *http.Request) (*estimateResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	request, err := api.DecodeBody[estimateMealRequest](r)
	if err != nil {
		return nil, 0, err
	}

	corrections := make([]ports.Correction, len(request.Corrections))
	for i, c := range request.Corrections {
		corrections[i] = ports.Correction{
			Item:       c.Item,
			Correction: c.Correction,
		}
	}

	estimate, err := h.useCase.Estimate(r.Context(), ports.EstimateInput{
		UserId:            userId,
		PhotoUrls:         request.PhotoUrls,
		Instructions:      request.Instructions,
		AssumeOnlyVisible: request.AssumeOnlyVisible,
		Corrections:       corrections,
	})
	if err != nil {
		return nil, 0, err
	}
	return estimateFromDomain(estimate), http.StatusOK, nil
}
