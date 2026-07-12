package use_case

import (
	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/authorization/ports"
	"github.com/ivan-ca97/life/internal/features/authorization/domain"
	authPorts "github.com/ivan-ca97/life/internal/features/authorization/ports"
	userPorts "github.com/ivan-ca97/life/internal/features/user/ports"
)

type shareUseCase struct {
	shareRepository authPorts.ShareRepository
	userService     userPorts.UserService
}

var _ ports.ShareUseCase = (*shareUseCase)(nil)

func NewShareUseCase(shareRepository authPorts.ShareRepository, userService userPorts.UserService) *shareUseCase {
	return &shareUseCase{
		shareRepository: shareRepository,
		userService:     userService,
	}
}

func (uc *shareUseCase) Create(ownerId uuid.UUID, granteeEmail, resourceType string, canWrite bool) (*domain.Share, error) {
	if !domain.ValidResourceTypes[resourceType] {
		return nil, cerr.NewBadRequestError("invalid resource_type")
	}

	grantee, err := uc.userService.GetByEmail(granteeEmail)
	if err != nil {
		return nil, cerr.NewBadRequestError("grantee user not found")
	}

	if grantee.Id == ownerId {
		return nil, cerr.NewBadRequestError("cannot share with yourself")
	}

	share := &domain.Share{
		Id:           uuid.New(),
		OwnerId:      ownerId,
		GranteeId:    grantee.Id,
		ResourceType: resourceType,
		CanWrite:     canWrite,
	}

	err = uc.shareRepository.Create(share)
	if err != nil {
		return nil, err
	}

	return share, nil
}

func (uc *shareUseCase) ListByOwner(ownerId uuid.UUID) ([]domain.Share, error) {
	shares, err := uc.shareRepository.ListByOwner(ownerId)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (uc *shareUseCase) ListByGrantee(granteeId uuid.UUID) ([]domain.Share, error) {
	shares, err := uc.shareRepository.ListByGrantee(granteeId)
	if err != nil {
		return nil, err
	}
	return shares, nil
}

func (uc *shareUseCase) Update(id, ownerId uuid.UUID, canWrite bool) (*domain.Share, error) {
	return uc.shareRepository.Update(id, ownerId, canWrite)
}

func (uc *shareUseCase) Delete(id, ownerId uuid.UUID) error {
	err := uc.shareRepository.Delete(id, ownerId)
	if err != nil {
		return err
	}
	return nil
}
