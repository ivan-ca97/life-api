package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/authentication/ports"
)

type sessionMiddleware struct {
	authenticationService ports.AuthenticationService
	errorHandler          http_errors.HttpErrorHandler
}

var _ api.Middleware = (*sessionMiddleware)(nil)

func NewSessionMiddleware(authenticationService ports.AuthenticationService, errorHandler http_errors.HttpErrorHandler) *sessionMiddleware {
	return &sessionMiddleware{
		authenticationService: authenticationService,
		errorHandler:          errorHandler,
	}
}

func (m *sessionMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r)
		if err != nil {
			m.errorHandler.Report(r, auth.ErrNoActor)
			return
		}

		session, err := m.authenticationService.Validate(token)
		if err != nil {
			m.errorHandler.Report(r, err)
			return
		}

		ctx := auth.WithActor(r.Context(), session.UserId)
		ctx = auth.WithSession(ctx, session.Id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(r *http.Request) (uuid.UUID, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return uuid.UUID{}, auth.ErrNoActor
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return uuid.UUID{}, auth.ErrNoActor
	}
	token, err := uuid.Parse(parts[1])
	if err != nil {
		return uuid.UUID{}, auth.ErrNoActor
	}
	return token, nil
}
