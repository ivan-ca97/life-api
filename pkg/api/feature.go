package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Feature interface {
	PublicRoutes(r chi.Router)
	ProtectedRoutes(r chi.Router)
}

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

// NoResponse is a sentinel type for handlers that return no body (204 No Content).
type NoResponse = *struct{}
