package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
)

type AuthorizedWeightEntryService interface {
	Create(ctx context.Context, ownerId uuid.UUID, params CreateParams) (*domain.WeightEntry, error)
	GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.WeightEntry, error)
	List(ctx context.Context, ownerId uuid.UUID, params ListParams) (types.Page[domain.WeightEntry], error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params UpdateParams) (*domain.WeightEntry, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
}
