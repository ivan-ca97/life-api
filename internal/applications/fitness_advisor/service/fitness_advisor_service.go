package service

import (
	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/domain"
	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/ports"
)

type fitnessAdvisorService struct {
	weightLookup ports.WeightLookup
}

var _ ports.FitnessAdvisorService = (*fitnessAdvisorService)(nil)

func NewFitnessAdvisorService(weightLookup ports.WeightLookup) *fitnessAdvisorService {
	return &fitnessAdvisorService{
		weightLookup: weightLookup,
	}
}

func (s *fitnessAdvisorService) EstimateCalories(userId uuid.UUID, request domain.EstimateRequest) (*domain.EstimateResult, error) {
	weightKg, err := s.weightLookup.LatestWeightKg(userId)
	if err != nil {
		return nil, err
	}
	if weightKg == nil {
		return nil, domain.ErrNoWeightData
	}

	calories, err := estimateByType(request, *weightKg)
	if err != nil {
		return nil, err
	}

	result := &domain.EstimateResult{
		Type:              request.Type,
		Value:             request.Value,
		EstimatedCalories: calories,
		WeightKg:          *weightKg,
	}
	return result, nil
}

func estimateByType(request domain.EstimateRequest, weightKg float64) (float64, error) {
	switch request.Type {
	case domain.ActivityTypeSteps:
		// MET=3.0 at cadence 100 steps/min (ACSM standard for moderate ambient walking)
		return request.Value * 0.0005 * weightKg, nil
	default:
		return 0, cerr.NewBadRequestError("unsupported activity type: " + string(request.Type))
	}
}
