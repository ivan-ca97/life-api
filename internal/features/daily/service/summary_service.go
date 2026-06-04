package service

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type summaryService struct {
	repository ports.SummaryRepository
}

var _ ports.SummaryService = (*summaryService)(nil)

func NewSummaryService(repository ports.SummaryRepository) *summaryService {
	return &summaryService{
		repository: repository,
	}
}

func (s *summaryService) GetSummary(userId uuid.UUID, date time.Time) (*domain.DailySummary, error) {
	summary, err := s.repository.GetDailySummary(userId, date)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (s *summaryService) GetSummaryRange(userId uuid.UUID, from, to time.Time) ([]domain.DailySummary, error) {
	summaries, err := s.repository.GetDailySummaryRange(userId, from, to)
	if err != nil {
		return nil, err
	}
	return summaries, nil
}
