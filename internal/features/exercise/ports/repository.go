package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
)

type ExerciseRepository interface {
	Create(exercise *domain.Exercise) error
	FindById(id, userId uuid.UUID) (*domain.Exercise, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.Exercise], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.Exercise, error)
	Delete(id, userId uuid.UUID) error
	ExistsByDateAndName(userId uuid.UUID, date time.Time, name string) (bool, error)
	FindByDateAndName(userId uuid.UUID, date time.Time, name string) (*domain.Exercise, error)
	ExistsByExternalId(userId uuid.UUID, externalId string) (bool, error)
}
