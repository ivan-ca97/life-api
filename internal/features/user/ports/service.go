package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type UpdateParams struct {
	Email     *string
	Password  *string
	HeightCm  *int
	BirthDate *time.Time
	Sex       *string
}

type UserService interface {
	Create(email, password string) (*domain.User, error)
	GetById(id uuid.UUID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	List(params types.PaginationParams) (types.Page[domain.User], error)
	Update(id uuid.UUID, params UpdateParams) (*domain.User, error)
	Deactivate(id uuid.UUID) error
}
