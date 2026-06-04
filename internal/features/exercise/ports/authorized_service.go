package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
)

type AuthorizedExerciseService interface {
	Create(ctx context.Context, params CreateParams) (*domain.Exercise, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Exercise, error)
	List(ctx context.Context, params ListParams) (types.Page[domain.Exercise], error)
	Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*domain.Exercise, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
