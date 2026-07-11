package ai_usage

import (
	"github.com/go-chi/chi/v5"

	"github.com/ivan-ca97/life/pkg/api/endpoint"
)

// Routes registers the AI usage endpoints. They live at the authenticated group
// level (not under /users/{userId}); admin endpoints are gated by the admin role
// inside the service, "me" endpoints act on the actor from context.
func (f *aiUsageFeature) Routes(r chi.Router) {
	r.Get("/ai/me/usage", endpoint.JSON(f.errorHandler, f.handler.GetMyUsage))
	r.Put("/ai/me/self-limit", endpoint.JSON(f.errorHandler, f.handler.SetMySelfLimit))

	r.Get("/ai/admin/tiers", endpoint.JSON(f.errorHandler, f.handler.ListTiers))
	r.Post("/ai/admin/tiers", endpoint.JSON(f.errorHandler, f.handler.CreateTier))
	r.Patch("/ai/admin/tiers/{tierId}", endpoint.JSON(f.errorHandler, f.handler.UpdateTier))
	r.Delete("/ai/admin/tiers/{tierId}", endpoint.JSON(f.errorHandler, f.handler.DeleteTier))
	r.Put("/ai/admin/users/{userId}/tier", endpoint.JSON(f.errorHandler, f.handler.AssignUserTier))
	r.Get("/ai/admin/users/{userId}/usage", endpoint.JSON(f.errorHandler, f.handler.GetUserUsage))
	r.Get("/ai/admin/interactions", endpoint.JSON(f.errorHandler, f.handler.ListInteractions))
}
