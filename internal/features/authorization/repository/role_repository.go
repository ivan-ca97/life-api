package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/authorization/ports"
)

type roleRepository struct {
	db *gorm.DB
}

var _ ports.RoleRepository = (*roleRepository)(nil)

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) UserHasPermission(userId uuid.UUID, permission string) (bool, error) {
	var count int64
	err := r.db.
		Table("role_permissions").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND role_permissions.permission = ?", userId, permission).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking user permission", err)
	}
	return count > 0, nil
}

func (r *roleRepository) UserHasRole(userId uuid.UUID, roleName string) (bool, error) {
	var count int64
	err := r.db.
		Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.name = ?", userId, roleName).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking user role", err)
	}
	return count > 0, nil
}

func (r *roleRepository) AssignRoleByName(userId uuid.UUID, roleName string) error {
	var rl role
	err := r.db.Where("name = ?", roleName).First(&rl).Error
	if err != nil {
		return cerr.NewInternalError("finding role by name", err)
	}
	ur := userRole{UserId: userId, RoleId: rl.Id}
	err = r.db.Create(&ur).Error
	if err != nil {
		return cerr.NewInternalError("assigning role to user", err)
	}
	return nil
}
