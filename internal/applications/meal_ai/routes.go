package meal_ai

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

// ProtectedRoutes registers the estimation endpoint. It is owner-scoped, so it
// lives under /users/{userId}: full path POST /users/{userId}/ai/meals/estimate.
func (a *MealAIApplication) ProtectedRoutes(r chi.Router) {
	r.Post("/ai/meals/estimate", endpoint.JSON(a.errorHandler, a.handler.Estimate))
}
