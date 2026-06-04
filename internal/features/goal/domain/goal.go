package domain

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	Id                   uuid.UUID
	UserId               uuid.UUID
	DailyCalories        *float64
	DailyProteinGrams    *float64
	DailyCarbsGrams      *float64
	DailyFatGrams        *float64
	DailyFiberGrams      *float64
	DailySteps           *int
	DailyExerciseMinutes *int
	TargetWeightKg       *float64
	StartedAt            time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
