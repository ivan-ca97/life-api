package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type UpsertParams struct {
	DailyCalories        *float64
	DailyProteinGrams    *float64
	DailyCarbsGrams      *float64
	DailyFatGrams        *float64
	DailyFiberGrams      *float64
	DailySteps           *int
	DailyExerciseMinutes *int
	TargetWeightKg       *float64
	StartedAt            *time.Time
}

type GoalService interface {
	GetByUserId(userId uuid.UUID) (*domain.Goal, error)
	Upsert(userId uuid.UUID, params UpsertParams) (*domain.Goal, error)
}
