package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
	"github.com/ivan-ca97/life/internal/features/goal/ports"
)

type goalService struct {
	repository ports.GoalRepository
}

var _ ports.GoalService = (*goalService)(nil)

func NewGoalService(repository ports.GoalRepository) *goalService {
	return &goalService{
		repository: repository,
	}
}

func (s *goalService) GetByUserId(userId uuid.UUID) (*domain.Goal, error) {
	goal, err := s.repository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *goalService) Upsert(userId uuid.UUID, params ports.UpsertParams) (*domain.Goal, error) {
	existing, err := s.repository.FindByUserId(userId)
	if err != nil && !errors.Is(err, domain.ErrGoalNotFound) {
		return nil, err
	}

	if existing != nil {
		if params.DailyCalories != nil {
			existing.DailyCalories = params.DailyCalories
		}
		if params.DailyProteinGrams != nil {
			existing.DailyProteinGrams = params.DailyProteinGrams
		}
		if params.DailyCarbsGrams != nil {
			existing.DailyCarbsGrams = params.DailyCarbsGrams
		}
		if params.DailyFatGrams != nil {
			existing.DailyFatGrams = params.DailyFatGrams
		}
		if params.DailyFiberGrams != nil {
			existing.DailyFiberGrams = params.DailyFiberGrams
		}
		if params.DailySteps != nil {
			existing.DailySteps = params.DailySteps
		}
		if params.DailyExerciseMinutes != nil {
			existing.DailyExerciseMinutes = params.DailyExerciseMinutes
		}
		if params.TargetWeightKg != nil {
			existing.TargetWeightKg = params.TargetWeightKg
		}
		if params.StartedAt != nil {
			existing.StartedAt = *params.StartedAt
		}
		goal, err := s.repository.Upsert(existing)
		if err != nil {
			return nil, err
		}
		return goal, nil
	}

	startedAt := time.Now()
	if params.StartedAt != nil {
		startedAt = *params.StartedAt
	}
	newGoal := &domain.Goal{
		Id:                   uuid.New(),
		UserId:               userId,
		DailyCalories:        params.DailyCalories,
		DailyProteinGrams:    params.DailyProteinGrams,
		DailyCarbsGrams:      params.DailyCarbsGrams,
		DailyFatGrams:        params.DailyFatGrams,
		DailyFiberGrams:      params.DailyFiberGrams,
		DailySteps:           params.DailySteps,
		DailyExerciseMinutes: params.DailyExerciseMinutes,
		TargetWeightKg:       params.TargetWeightKg,
		StartedAt:            startedAt,
	}
	goal, err := s.repository.Upsert(newGoal)
	if err != nil {
		return nil, err
	}
	return goal, nil
}
