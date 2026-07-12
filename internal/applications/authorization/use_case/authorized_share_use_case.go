package use_case

import (
	"context"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/applications/authorization/ports"
	"github.com/ivan-ca97/life/internal/features/authorization/domain"
	"github.com/ivan-ca97/life/internal/permissions"
)

type authorizedShareUseCase struct {
	base       ports.ShareUseCase
	authorizer auth.AuthorizationService
}

var _ ports.AuthorizedShareUseCase = (*authorizedShareUseCase)(nil)

func NewAuthorizedShareUseCase(base ports.ShareUseCase, authorizer auth.AuthorizationService) *authorizedShareUseCase {
	return &authorizedShareUseCase{
		base:       base,
		authorizer: authorizer,
	}
}

func (uc *authorizedShareUseCase) Create(ctx context.Context, ownerId uuid.UUID, granteeEmail, resourceType string, canWrite bool) (*domain.Share, error) {
	err := uc.authorizer.Authorize(ctx, ownerId, permissions.SharesCreate)
	if err != nil {
		return nil, err
	}
	share, err := uc.base.Create(ownerId, granteeEmail, resourceType, canWrite)
	if err != nil {
		return nil, err
	}
	return share, nil
}

func (uc *authorizedShareUseCase) ListOwned(ctx context.Context, ownerId uuid.UUID) ([]domain.Share, error) {
	err := uc.authorizer.Authorize(ctx, ownerId, permissions.SharesRead)
	if err != nil {
		return nil, err
	}
	shares, err := uc.base.ListByOwner(ownerId)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (uc *authorizedShareUseCase) ListReceived(ctx context.Context, ownerId uuid.UUID) ([]domain.Share, error) {
	err := uc.authorizer.Authorize(ctx, ownerId, permissions.SharesRead)
	if err != nil {
		return nil, err
	}
	shares, err := uc.base.ListByGrantee(ownerId)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (uc *authorizedShareUseCase) Update(ctx context.Context, ownerId uuid.UUID, id uuid.UUID, canWrite bool) (*domain.Share, error) {
	if err := uc.authorizer.Authorize(ctx, ownerId, permissions.SharesUpdate); err != nil {
		return nil, err
	}
	return uc.base.Update(id, ownerId, canWrite)
}

func (uc *authorizedShareUseCase) Delete(ctx context.Context, ownerId uuid.UUID, id uuid.UUID) error {
	err := uc.authorizer.Authorize(ctx, ownerId, permissions.SharesDelete)
	if err != nil {
		return err
	}
	err = uc.base.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return nil
}
