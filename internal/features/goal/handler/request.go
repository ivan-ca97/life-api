package handler

type upsertGoalRequest struct {
	DailyCalories        *float64 `json:"daily_calories,omitempty"`
	DailyProteinGrams    *float64 `json:"daily_protein_grams,omitempty"`
	DailyCarbsGrams      *float64 `json:"daily_carbs_grams,omitempty"`
	DailyFatGrams        *float64 `json:"daily_fat_grams,omitempty"`
	DailyFiberGrams      *float64 `json:"daily_fiber_grams,omitempty"`
	DailySteps           *int     `json:"daily_steps,omitempty"`
	DailyExerciseMinutes *int     `json:"daily_exercise_minutes,omitempty"`
	TargetWeightKg       *float64 `json:"target_weight_kg,omitempty"`
	StartedAt            *string  `json:"started_at,omitempty"`
}
