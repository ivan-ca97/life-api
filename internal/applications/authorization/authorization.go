package authorization

import (
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	authPorts "github.com/ivan-ca97/life/internal/features/authorization/ports"
	userPorts "github.com/ivan-ca97/life/internal/features/user/ports"

	"github.com/ivan-ca97/life/internal/applications/authorization/handler"
	"github.com/ivan-ca97/life/internal/applications/authorization/use_case"
)

type authorizationApplication struct {
	shareHandler handler.ShareHandler
	errorHandler http_errors.HttpErrorHandler
}

func NewAuthorizationApplication(
	shareRepository authPorts.ShareRepository,
	authorizer auth.AuthorizationService,
	userService userPorts.UserService,
	errorHandler http_errors.HttpErrorHandler,
) *authorizationApplication {
	shareUseCase := use_case.NewShareUseCase(shareRepository, userService)
	authorizedShareUseCase := use_case.NewAuthorizedShareUseCase(shareUseCase, authorizer)
	shareHandler := handler.NewShareHandler(authorizedShareUseCase)

	return &authorizationApplication{
		shareHandler: shareHandler,
		errorHandler: errorHandler,
	}
}
