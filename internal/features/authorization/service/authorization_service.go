package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/authorization/ports"
)

type authorizationService struct {
	roleRepo  ports.RoleRepository
	shareRepo ports.ShareRepository
}

var _ auth.AuthorizationService = (*authorizationService)(nil)

func NewAuthorizationService(roleRepo ports.RoleRepository, shareRepo ports.ShareRepository) *authorizationService {
	return &authorizationService{
		roleRepo:  roleRepo,
		shareRepo: shareRepo,
	}
}

func (s *authorizationService) AuthorizeAdmin(ctx context.Context) error {
	actorId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return auth.ErrNoActor
	}
	isAdmin, err := s.roleRepo.UserHasRole(actorId, "admin")
	if err != nil {
		return err
	}
	if !isAdmin {
		return auth.ErrForbidden
	}
	return nil
}

func (s *authorizationService) Authorize(ctx context.Context, ownerId uuid.UUID, permission string) error {
	actorId, err := auth.ActorFromContext(ctx)
	if err != nil {
		return auth.ErrNoActor
	}

	// Own data: check role permissions
	if actorId == ownerId {
		hasPermission, err := s.roleRepo.UserHasPermission(actorId, permission)
		if err != nil {
			return err
		}
		if !hasPermission {
			return auth.ErrForbidden
		}
		return nil
	}

	// Cross-user access: check admin role first
	isAdmin, err := s.roleRepo.UserHasRole(actorId, "admin")
	if err != nil {
		return err
	}
	if isAdmin {
		return nil
	}

	// Check shares
	resourceType, needsWrite := parsePermission(permission)
	hasAccess, err := s.shareRepo.HasAccess(ownerId, actorId, resourceType, needsWrite)
	if err != nil {
		return err
	}
	if !hasAccess {
		return auth.ErrForbidden
	}
	return nil
}

// parsePermission splits "meals:read" into ("meals", false) and "meals:create" into ("meals", true).
func parsePermission(permission string) (resourceType string, needsWrite bool) {
	parts := strings.SplitN(permission, ":", 2)
	if len(parts) != 2 {
		return permission, true
	}
	resourceType = parts[0]
	needsWrite = parts[1] != "read"
	return
}
