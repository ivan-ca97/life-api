package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authorization/domain"
)

type RoleRepository interface {
	UserHasPermission(userId uuid.UUID, permission string) (bool, error)
	UserHasRole(userId uuid.UUID, roleName string) (bool, error)
	AssignRoleByName(userId uuid.UUID, roleName string) error
}

type ShareRepository interface {
	HasAccess(ownerId, granteeId uuid.UUID, resourceType string, needsWrite bool) (bool, error)
	Create(share *domain.Share) error
	ListByOwner(ownerId uuid.UUID) ([]domain.Share, error)
	ListByGrantee(granteeId uuid.UUID) ([]domain.Share, error)
	Delete(id, ownerId uuid.UUID) error
}
