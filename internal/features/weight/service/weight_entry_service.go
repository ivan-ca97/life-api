package service

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
	"github.com/ivan-ca97/life/internal/features/weight/ports"
)

type weightEntryService struct {
	repository ports.WeightEntryRepository
}

var _ ports.WeightEntryService = (*weightEntryService)(nil)

func NewWeightEntryService(repository ports.WeightEntryRepository) *weightEntryService {
	return &weightEntryService{
		repository: repository,
	}
}

func (s *weightEntryService) Create(userId uuid.UUID, params ports.CreateParams) (*domain.WeightEntry, error) {
	entry := &domain.WeightEntry{
		Id:                uuid.New(),
		UserId:            userId,
		Date:              params.Date,
		WeightKg:          params.WeightKg,
		BodyFatPercentage: params.BodyFatPercentage,
		Notes:             params.Notes,
	}
	err := s.repository.Create(entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *weightEntryService) GetById(id, userId uuid.UUID) (*domain.WeightEntry, error) {
	entry, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *weightEntryService) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.WeightEntry], error) {
	page, err := s.repository.List(userId, params)
	if err != nil {
		return types.Page[domain.WeightEntry]{}, err
	}
	return page, nil
}

func (s *weightEntryService) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.WeightEntry, error) {
	entry, err := s.repository.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *weightEntryService) Delete(id, userId uuid.UUID) error {
	err := s.repository.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}
