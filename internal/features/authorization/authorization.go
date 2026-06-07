package authorization

import (
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/features/authorization/ports"
	"github.com/ivan-ca97/life/internal/features/authorization/repository"
	"github.com/ivan-ca97/life/internal/features/authorization/service"
)

type AuthorizationFeature struct {
	authorizationService auth.AuthorizationService
	roleRepository       ports.RoleRepository
	shareRepository      ports.ShareRepository
}

func NewAuthorizationFeature(db *gorm.DB) *AuthorizationFeature {
	roleRepository := repository.NewRoleRepository(db)
	shareRepository := repository.NewShareRepository(db)
	authorizationService := service.NewAuthorizationService(roleRepository, shareRepository)

	return &AuthorizationFeature{
		authorizationService: authorizationService,
		roleRepository:       roleRepository,
		shareRepository:      shareRepository,
	}
}

func (f *AuthorizationFeature) AuthorizationService() auth.AuthorizationService {
	return f.authorizationService
}

func (f *AuthorizationFeature) RoleRepository() ports.RoleRepository {
	return f.roleRepository
}

func (f *AuthorizationFeature) ShareRepository() ports.ShareRepository {
	return f.shareRepository
}
