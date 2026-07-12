package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/authorization/domain"
)

type ShareUseCase interface {
	Create(ownerId uuid.UUID, granteeEmail, resourceType string, canWrite bool) (*domain.Share, error)
	ListByOwner(ownerId uuid.UUID) ([]domain.Share, error)
	ListByGrantee(granteeId uuid.UUID) ([]domain.Share, error)
	Update(id, ownerId uuid.UUID, canWrite bool) (*domain.Share, error)
	Delete(id, ownerId uuid.UUID) error
}

type AuthorizedShareUseCase interface {
	Create(ctx context.Context, ownerId uuid.UUID, granteeEmail, resourceType string, canWrite bool) (*domain.Share, error)
	ListOwned(ctx context.Context, ownerId uuid.UUID) ([]domain.Share, error)
	ListReceived(ctx context.Context, ownerId uuid.UUID) ([]domain.Share, error)
	Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, canWrite bool) (*domain.Share, error)
	Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error
}
