package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/user/domain"
)

type ProfilePhotoRepository interface {
	Create(photo *domain.ProfilePhoto) error
	ListByUserId(userId uuid.UUID, params types.PaginationParams) (types.Page[domain.ProfilePhoto], error)
}
