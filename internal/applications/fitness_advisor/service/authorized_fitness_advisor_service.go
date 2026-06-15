package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/domain"
	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/ports"
)

type authorizedFitnessAdvisorService struct {
	base       ports.FitnessAdvisorService
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedFitnessAdvisorService = (*authorizedFitnessAdvisorService)(nil)

func NewAuthorizedFitnessAdvisorService(base ports.FitnessAdvisorService, authorizer auth.AuthorizationService) *authorizedFitnessAdvisorService {
	return &authorizedFitnessAdvisorService{
		base:       base,
		authorizer: authorizer,
	}
}

func (s *authorizedFitnessAdvisorService) EstimateCalories(ctx context.Context, userId uuid.UUID, request domain.EstimateRequest) (*domain.EstimateResult, error) {
	err := s.authorizer.Authorize(ctx, userId, "read")
	if err != nil {
		return nil, err
	}

	result, err := s.base.EstimateCalories(userId, request)
	if err != nil {
		return nil, err
	}
	return result, nil
}
