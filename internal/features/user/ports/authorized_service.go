package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type AuthorizedUserService interface {
	Create(ctx context.Context, email, password string) (*domain.User, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, params types.PaginationParams) (types.Page[domain.User], error)
	Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*domain.User, error)
	Deactivate(ctx context.Context, id uuid.UUID) error
}
