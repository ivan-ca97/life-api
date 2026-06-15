package ports

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
)

type WeightEntryRepository interface {
	Create(entry *domain.WeightEntry) error
	FindById(id, userId uuid.UUID) (*domain.WeightEntry, error)
	LatestByUserId(userId uuid.UUID) (*domain.WeightEntry, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.WeightEntry], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.WeightEntry, error)
	Delete(id, userId uuid.UUID) error
	ExistsByExternalId(userId uuid.UUID, externalId string) (bool, error)
}
