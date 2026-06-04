package auth

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/api/http_errors"

	"github.com/ivan-ca97/life/internal/features/auth/handler"
	"github.com/ivan-ca97/life/internal/features/auth/middleware"
	"github.com/ivan-ca97/life/internal/features/auth/ports"
	"github.com/ivan-ca97/life/internal/features/auth/repository"
	"github.com/ivan-ca97/life/internal/features/auth/service"

	user_ports "github.com/ivan-ca97/life/internal/features/user/ports"
)

type authFeature struct {
	service      ports.AuthService
	middleware   api.Middleware
	authHandler  handler.AuthHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewAuthFeature(db *gorm.DB, userService user_ports.UserService, errorHandler http_errors.HttpErrorHandler) *authFeature {
	sessionRepository := repository.NewSessionRepository(db)
	authService := service.NewAuthService(sessionRepository, userService)
	authHandler := handler.NewAuthHandler(authService)
	sessionMiddleware := middleware.NewSessionMiddleware(authService, errorHandler)

	return &authFeature{
		service:      authService,
		middleware:   sessionMiddleware,
		authHandler:  authHandler,
		errorHandler: errorHandler,
	}
}

func (f *authFeature) Service() ports.AuthService {
	return f.service
}

func (f *authFeature) Middleware() api.Middleware {
	return f.middleware
}
