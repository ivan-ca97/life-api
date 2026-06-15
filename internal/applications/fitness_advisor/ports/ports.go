package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/domain"
)

type FitnessAdvisorService interface {
	EstimateCalories(userId uuid.UUID, request domain.EstimateRequest) (*domain.EstimateResult, error)
}

type AuthorizedFitnessAdvisorService interface {
	EstimateCalories(ctx context.Context, userId uuid.UUID, request domain.EstimateRequest) (*domain.EstimateResult, error)
}

type WeightLookup interface {
	LatestWeightKg(userId uuid.UUID) (*float64, error)
}
