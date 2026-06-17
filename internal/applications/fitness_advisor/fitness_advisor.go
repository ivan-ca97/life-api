package fitness_advisor

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/handler"
	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/ports"
	"github.com/ivan-ca97/life/internal/applications/fitness_advisor/service"
	weightPorts "github.com/ivan-ca97/life/internal/features/weight/ports"
)

type FitnessAdvisorApplication struct {
	handler      handler.FitnessAdvisorHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewFitnessAdvisorApplication(
	weightRepository weightPorts.WeightEntryRepository,
	authorizer auth.AuthorizationService,
	errorHandler http_errors.HttpErrorHandler,
) *FitnessAdvisorApplication {
	lookup := newWeightRepositoryLookup(weightRepository)
	advisorService := service.NewFitnessAdvisorService(lookup)
	authorizedService := service.NewAuthorizedFitnessAdvisorService(advisorService, authorizer)
	advisorHandler := handler.NewFitnessAdvisorHandler(authorizedService)

	application := &FitnessAdvisorApplication{
		handler:      advisorHandler,
		errorHandler: errorHandler,
	}
	return application
}

type weightRepositoryLookup struct {
	repository weightPorts.WeightEntryRepository
}

var _ ports.WeightLookup = (*weightRepositoryLookup)(nil)

func newWeightRepositoryLookup(repository weightPorts.WeightEntryRepository) *weightRepositoryLookup {
	lookup := &weightRepositoryLookup{
		repository: repository,
	}
	return lookup
}

func (w *weightRepositoryLookup) LatestWeightKg(userId uuid.UUID) (*float64, error) {
	entry, err := w.repository.LatestByUserId(userId)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	weight := entry.WeightKg
	return &weight, nil
}
