package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
)

type AuthorizedExerciseService interface {
	Create(ctx context.Context, ownerId uuid.UUID, params CreateParams) (*domain.Exercise, error)
	GetById(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) (*domain.Exercise, error)
	List(ctx context.Context, ownerId uuid.UUID, params ListParams) (types.Page[domain.Exercise], error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, params UpdateParams) (*domain.Exercise, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
}
