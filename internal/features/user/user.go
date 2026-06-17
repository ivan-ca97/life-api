package user

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/user/handler"
	"github.com/ivan-ca97/life/internal/features/user/ports"
	"github.com/ivan-ca97/life/internal/features/user/repository"
	"github.com/ivan-ca97/life/internal/features/user/service"
)

type userFeature struct {
	service           ports.UserService
	authorizedService ports.AuthorizedUserService
	userHandler       handler.UserHandler
	errorHandler      http_errors.HttpErrorHandler
}

func NewUserFeature(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *userFeature {
	userRepository := repository.NewUserRepository(db)
	profilePhotoRepository := repository.NewProfilePhotoRepository(db)
	userService := service.NewUserService(userRepository, profilePhotoRepository)
	authorizedService := service.NewAuthorizedUserService(userService, authorizer)
	userHandler := handler.NewUserHandler(authorizedService)

	return &userFeature{
		service:           userService,
		authorizedService: authorizedService,
		userHandler:       userHandler,
		errorHandler:      errorHandler,
	}
}

func (f *userFeature) Service() ports.UserService {
	return f.service
}

func (f *userFeature) AuthorizedService() ports.AuthorizedUserService {
	return f.authorizedService
}
