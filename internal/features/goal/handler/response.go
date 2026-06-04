package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
)

type goalResponse struct {
	Id                   uuid.UUID `json:"id"`
	DailyCalories        *float64  `json:"daily_calories,omitempty"`
	DailyProteinGrams    *float64  `json:"daily_protein_grams,omitempty"`
	DailyCarbsGrams      *float64  `json:"daily_carbs_grams,omitempty"`
	DailyFatGrams        *float64  `json:"daily_fat_grams,omitempty"`
	DailyFiberGrams      *float64  `json:"daily_fiber_grams,omitempty"`
	DailySteps           *int      `json:"daily_steps,omitempty"`
	DailyExerciseMinutes *int      `json:"daily_exercise_minutes,omitempty"`
	TargetWeightKg       *float64  `json:"target_weight_kg,omitempty"`
	StartedAt            time.Time `json:"started_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func goalFromDomain(g *domain.Goal) *goalResponse {
	return &goalResponse{
		Id:                   g.Id,
		DailyCalories:        g.DailyCalories,
		DailyProteinGrams:    g.DailyProteinGrams,
		DailyCarbsGrams:      g.DailyCarbsGrams,
		DailyFatGrams:        g.DailyFatGrams,
		DailyFiberGrams:      g.DailyFiberGrams,
		DailySteps:           g.DailySteps,
		DailyExerciseMinutes: g.DailyExerciseMinutes,
		TargetWeightKg:       g.TargetWeightKg,
		StartedAt:            g.StartedAt,
		CreatedAt:            g.CreatedAt,
		UpdatedAt:            g.UpdatedAt,
	}
}
