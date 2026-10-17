package nutrition

import (
	"github.com/go-chi/chi/v5"
)

func (f *nutritionFeature) PublicRoutes(r chi.Router) {
	f.foods.PublicRoutes(r)
	f.meals.PublicRoutes(r)
}

// GlobalRoutes are the routes that are not scoped to a user.
func (f *nutritionFeature) GlobalRoutes(r chi.Router) {
	f.foods.GlobalRoutes(r)
}

func (f *nutritionFeature) ProtectedRoutes(r chi.Router) {
	f.foods.ProtectedRoutes(r)
	f.meals.ProtectedRoutes(r)
}
