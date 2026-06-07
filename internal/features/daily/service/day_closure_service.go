package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type dayClosureService struct {
	repository ports.DayClosureRepository
}

var _ ports.DayClosureService = (*dayClosureService)(nil)
var _ dayclosure.DayClosureChecker = (*dayClosureService)(nil)

func NewDayClosureService(repository ports.DayClosureRepository) *dayClosureService {
	return &dayClosureService{repository: repository}
}

func (s *dayClosureService) Close(userId uuid.UUID, date time.Time) error {
	err := s.repository.Close(userId, date)
	if err != nil {
		return err
	}
	return nil
}

func (s *dayClosureService) Open(userId uuid.UUID, date time.Time) error {
	err := s.repository.Open(userId, date)
	if err != nil {
		return err
	}
	return nil
}

func (s *dayClosureService) IsClosed(userId uuid.UUID, date time.Time) (bool, error) {
	closed, err := s.repository.IsClosed(userId, date)
	if err != nil {
		return false, err
	}
	return closed, nil
}

func (s *dayClosureService) GetClosedDates(userId uuid.UUID, from, to time.Time) (map[string]bool, error) {
	dates, err := s.repository.GetClosedDates(userId, from, to)
	if err != nil {
		return nil, err
	}
	return dates, nil
}
