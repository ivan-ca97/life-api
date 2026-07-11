package service

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/dayclosure"
	"github.com/ivan-ca97/life/pkg/types"
	"github.com/ivan-ca97/life/pkg/validate"

	"github.com/ivan-ca97/life/internal/features/weight/domain"
	"github.com/ivan-ca97/life/internal/features/weight/ports"
)

type weightEntryService struct {
	repository     ports.WeightEntryRepository
	closureChecker dayclosure.DayClosureChecker
}

var _ ports.WeightEntryService = (*weightEntryService)(nil)

func NewWeightEntryService(repository ports.WeightEntryRepository, closureChecker dayclosure.DayClosureChecker) *weightEntryService {
	return &weightEntryService{
		repository:     repository,
		closureChecker: closureChecker,
	}
}

func (s *weightEntryService) Create(userId uuid.UUID, params ports.CreateParams) (*domain.WeightEntry, error) {
	err := validate.InRange(params.WeightKg, 30, 500, "weight_kg")
	if err != nil {
		return nil, err
	}
	err = validate.InRangePtr(params.BodyFatPercentage, 0, 100, "body_fat_percentage")
	if err != nil {
		return nil, err
	}
	closed, err := s.closureChecker.IsClosed(userId, params.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	entry := &domain.WeightEntry{
		Id:                uuid.New(),
		UserId:            userId,
		Date:              params.Date,
		WeightKg:          params.WeightKg,
		BodyFatPercentage: params.BodyFatPercentage,
		Notes:             params.Notes,
		ExternalId:        params.ExternalId,
	}
	err = s.repository.Create(entry)
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
	err := validate.InRangePtr(params.WeightKg, 30, 500, "weight_kg")
	if err != nil {
		return nil, err
	}
	err = validate.InRangePtr(params.BodyFatPercentage, 0, 100, "body_fat_percentage")
	if err != nil {
		return nil, err
	}
	entry, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	closed, err := s.closureChecker.IsClosed(userId, entry.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	updated, err := s.repository.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *weightEntryService) Delete(id, userId uuid.UUID) error {
	entry, err := s.repository.FindById(id, userId)
	if err != nil {
		return err
	}
	closed, err := s.closureChecker.IsClosed(userId, entry.Date)
	if err != nil {
		return err
	}
	if closed {
		return dayclosure.ErrDayClosed
	}

	err = s.repository.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}
