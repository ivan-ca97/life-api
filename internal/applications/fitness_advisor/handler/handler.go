package handler

import (
	"errors"
	"net/http"

	"github.com/ivan-ca97/life/pkg/api"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/domain"
	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/ports"
)

type FitnessAdvisorHandler interface {
	EstimateCalories(r *http.Request) (*estimateResponse, int, error)
}

type fitnessAdvisorHandler struct {
	service ports.AuthorizedFitnessAdvisorService
}

var _ FitnessAdvisorHandler = (*fitnessAdvisorHandler)(nil)

func NewFitnessAdvisorHandler(service ports.AuthorizedFitnessAdvisorService) *fitnessAdvisorHandler {
	handler := &fitnessAdvisorHandler{
		service: service,
	}
	return handler
}

func (h *fitnessAdvisorHandler) EstimateCalories(r *http.Request) (*estimateResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}

	request, err := api.DecodeBody[estimateRequest](r)
	if err != nil {
		return nil, 0, err
	}

	if request.Value <= 0 {
		return nil, 0, cerr.NewBadRequestError("value must be positive")
	}

	estimateReq := domain.EstimateRequest{
		Type:  domain.ActivityType(request.Type),
		Value: request.Value,
	}

	result, err := h.service.EstimateCalories(r.Context(), userId, estimateReq)
	if err != nil {
		if errors.Is(err, domain.ErrNoWeightData) {
			return nil, 0, cerr.NewBadRequestError(domain.ErrNoWeightData.Error())
		}
		if errors.Is(err, domain.ErrUnsupportedActivityType) {
			return nil, 0, cerr.NewBadRequestError("unsupported activity type: " + string(estimateReq.Type))
		}
		return nil, 0, err
	}

	response := estimateResponseFromDomain(result)
	return response, http.StatusOK, nil
}
