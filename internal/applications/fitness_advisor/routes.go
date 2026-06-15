package fitness_advisor

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

func (a *FitnessAdvisorApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/activity/calorie-estimate", endpoint.JSON(a.errorHandler, a.handler.EstimateCalories))
}
