package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type AuthorizedUserService interface {
	Create(ctx context.Context, email, password string) (*domain.User, error)
	GetById(ctx context.Context, ownerId uuid.UUID) (*domain.User, error)
	List(ctx context.Context, params types.PaginationParams) (types.Page[domain.User], error)
	Update(ctx context.Context, ownerId uuid.UUID, params UpdateParams) (*domain.User, error)
	Deactivate(ctx context.Context, ownerId uuid.UUID) error
	AddProfilePhoto(ctx context.Context, userId uuid.UUID, url string) (*domain.ProfilePhoto, error)
	ListProfilePhotos(ctx context.Context, userId uuid.UUID, params types.PaginationParams) (types.Page[domain.ProfilePhoto], error)
}
