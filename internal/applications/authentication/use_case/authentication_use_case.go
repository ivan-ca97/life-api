package use_case

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	appDomain "github.com/ivan-ca97/life/internal/applications/authentication/domain"
	"github.com/ivan-ca97/life/internal/applications/authentication/ports"
	authenticationPorts "github.com/ivan-ca97/life/internal/features/authentication/ports"
	userDomain "github.com/ivan-ca97/life/internal/features/user/domain"
	userPorts "github.com/ivan-ca97/life/internal/features/user/ports"
)

type authenticationUseCase struct {
	authenticationService authenticationPorts.AuthenticationService
	userService           userPorts.UserService
	roleAssigner          ports.RoleAssigner
	googleVerifier        authenticationPorts.GoogleTokenVerifier
	googleClientId        string
}

var _ ports.AuthenticationUseCase = (*authenticationUseCase)(nil)

func NewAuthenticationUseCase(
	authenticationService authenticationPorts.AuthenticationService,
	userService userPorts.UserService,
	roleAssigner ports.RoleAssigner,
	googleVerifier authenticationPorts.GoogleTokenVerifier,
	googleClientId string,
) *authenticationUseCase {
	return &authenticationUseCase{
		authenticationService: authenticationService,
		userService:           userService,
		roleAssigner:          roleAssigner,
		googleVerifier:        googleVerifier,
		googleClientId:        googleClientId,
	}
}

func (uc *authenticationUseCase) Login(email, password string) (*ports.AuthenticationResult, error) {
	user, err := uc.userService.GetByEmail(email)
	if err != nil {
		return nil, appDomain.ErrInvalidCredentials
	}

	if !user.Active {
		return nil, userDomain.ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, appDomain.ErrInvalidCredentials
	}

	session, err := uc.authenticationService.CreateSession(user.Id)
	if err != nil {
		return nil, err
	}

	result := &ports.AuthenticationResult{
		UserId:    session.UserId,
		Token:     session.Id,
		ExpiresAt: session.ExpiresAt,
	}
	return result, nil
}

func (uc *authenticationUseCase) Register(email, password string) (*ports.AuthenticationResult, error) {
	user, err := uc.userService.Create(email, password)
	if err != nil {
		return nil, err
	}

	err = uc.roleAssigner.AssignRoleByName(user.Id, "user")
	if err != nil {
		return nil, err
	}

	session, err := uc.authenticationService.CreateSession(user.Id)
	if err != nil {
		return nil, err
	}

	result := &ports.AuthenticationResult{
		UserId:    user.Id,
		Token:     session.Id,
		ExpiresAt: session.ExpiresAt,
	}
	return result, nil
}

func (uc *authenticationUseCase) LoginWithGoogle(idToken string) (*ports.AuthenticationResult, error) {
	claims, err := uc.googleVerifier.Verify(idToken, uc.googleClientId)
	if err != nil {
		return nil, err
	}

	user, err := uc.userService.GetByEmail(claims.Email)
	if err != nil {
		if !errors.Is(err, userDomain.ErrUserNotFound) {
			return nil, err
		}

		user, err = uc.userService.CreateOAuth(claims.Email, claims.Subject)
		if err != nil {
			return nil, err
		}

		err = uc.roleAssigner.AssignRoleByName(user.Id, "user")
		if err != nil {
			return nil, err
		}
	}

	if !user.Active {
		return nil, userDomain.ErrUserInactive
	}

	session, err := uc.authenticationService.CreateSession(user.Id)
	if err != nil {
		return nil, err
	}

	result := &ports.AuthenticationResult{
		UserId:    session.UserId,
		Token:     session.Id,
		ExpiresAt: session.ExpiresAt,
	}
	return result, nil
}

func (uc *authenticationUseCase) Logout(sessionId uuid.UUID) error {
	err := uc.authenticationService.Logout(sessionId)
	if err != nil {
		return err
	}
	return nil
}
