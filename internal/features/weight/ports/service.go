package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
)

type CreateParams struct {
	Date              time.Time
	WeightKg          float64
	BodyFatPercentage *float64
	Notes             string
	ExternalId        *string
}

type UpdateParams struct {
	Date              *time.Time
	WeightKg          *float64
	BodyFatPercentage *float64
	Notes             *string
}

type ListParams struct {
	types.PaginationParams
	From *time.Time
	To   *time.Time
}

type WeightEntryService interface {
	Create(userId uuid.UUID, params CreateParams) (*domain.WeightEntry, error)
	GetById(id, userId uuid.UUID) (*domain.WeightEntry, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.WeightEntry], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.WeightEntry, error)
	Delete(id, userId uuid.UUID) error
}
