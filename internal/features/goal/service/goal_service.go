package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/validate"

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
	err := validate.NonNegativePtr(params.DailyCalories, "daily_calories")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.DailyProteinGrams, "daily_protein_grams")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.DailyCarbsGrams, "daily_carbs_grams")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.DailyFatGrams, "daily_fat_grams")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.DailyFiberGrams, "daily_fiber_grams")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.DailySteps, "daily_steps")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.DailyExerciseMinutes, "daily_exercise_minutes")
	if err != nil {
		return nil, err
	}
	err = validate.PositivePtr(params.TargetWeightKg, "target_weight_kg")
	if err != nil {
		return nil, err
	}
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

func (s *goalService) GetProgress(userId uuid.UUID, from, to time.Time) (*domain.GoalProgress, error) {
	return s.repository.GetProgress(userId, from, to)
}
