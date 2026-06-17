package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindById(id uuid.UUID) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	List(params types.PaginationParams) (types.Page[domain.User], error)
	Update(id uuid.UUID, params UpdateParams) (*domain.User, error)
	UpdatePhotoUrl(id uuid.UUID, url string) error
	Deactivate(id uuid.UUID) error
	EmailExists(email string) (bool, error)
}
