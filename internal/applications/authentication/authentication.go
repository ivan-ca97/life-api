package authentication

import (
	"github.com/ivan-ca97/life/pkg/api/http_errors"

	authenticationPorts "github.com/ivan-ca97/life/internal/features/authentication/ports"
	userPorts "github.com/ivan-ca97/life/internal/features/user/ports"

	"github.com/ivan-ca97/life/internal/applications/authentication/handler"
	"github.com/ivan-ca97/life/internal/applications/authentication/ports"
	"github.com/ivan-ca97/life/internal/applications/authentication/use_case"
)

type authenticationApplication struct {
	authenticationHandler handler.AuthenticationHandler
	errorHandler          http_errors.HttpErrorHandler
}

func NewAuthenticationApplication(
	authenticationService authenticationPorts.AuthenticationService,
	userService userPorts.UserService,
	roleAssigner ports.RoleAssigner,
	googleVerifier authenticationPorts.GoogleTokenVerifier,
	googleClientId string,
	errorHandler http_errors.HttpErrorHandler,
) *authenticationApplication {
	authenticationUseCase := use_case.NewAuthenticationUseCase(
		authenticationService,
		userService,
		roleAssigner,
		googleVerifier,
		googleClientId,
	)
	authenticationHandler := handler.NewAuthenticationHandler(authenticationUseCase)

	return &authenticationApplication{
		authenticationHandler: authenticationHandler,
		errorHandler:          errorHandler,
	}
}
