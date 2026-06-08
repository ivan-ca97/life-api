package authentication

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/api/http_errors"

	"github.com/ivan-ca97/life/internal/features/authentication/middleware"
	"github.com/ivan-ca97/life/internal/features/authentication/ports"
	"github.com/ivan-ca97/life/internal/features/authentication/repository"
	"github.com/ivan-ca97/life/internal/features/authentication/service"
)

type AuthenticationFeature struct {
	authenticationService ports.AuthenticationService
	middleware            api.Middleware
	googleVerifier        ports.GoogleTokenVerifier
}

func NewAuthenticationFeature(db *gorm.DB, errorHandler http_errors.HttpErrorHandler) *AuthenticationFeature {
	sessionRepository := repository.NewSessionRepository(db)
	authenticationService := service.NewAuthenticationService(sessionRepository)
	googleVerifier := service.NewGoogleTokenVerifier()
	sessionMiddleware := middleware.NewSessionMiddleware(authenticationService, errorHandler)

	return &AuthenticationFeature{
		authenticationService: authenticationService,
		middleware:            sessionMiddleware,
		googleVerifier:        googleVerifier,
	}
}

func (f *AuthenticationFeature) Service() ports.AuthenticationService {
	return f.authenticationService
}

func (f *AuthenticationFeature) Middleware() api.Middleware {
	return f.middleware
}

func (f *AuthenticationFeature) GoogleVerifier() ports.GoogleTokenVerifier {
	return f.googleVerifier
}
