package handler

import "github.com/ivan-ca97/life/internal/features/daily/domain"

type mealsSummaryResponse struct {
	TotalCalories     float64 `json:"total_calories"`
	TotalProteinGrams float64 `json:"total_protein_grams"`
	TotalCarbsGrams   float64 `json:"total_carbs_grams"`
	TotalFatGrams     float64 `json:"total_fat_grams"`
	TotalFiberGrams   float64 `json:"total_fiber_grams"`
	Count             int     `json:"count"`
}

type exerciseSummaryResponse struct {
	TotalCaloriesBurned  float64 `json:"total_calories_burned"`
	TotalSteps           int     `json:"total_steps"`
	TotalDurationSeconds int     `json:"total_duration_seconds"`
	TotalDistanceMeters  float64 `json:"total_distance_meters"`
	Count                int     `json:"count"`
}

type weightEntrySummaryResponse struct {
	WeightKg          float64  `json:"weight_kg"`
	BodyFatPercentage *float64 `json:"body_fat_percentage,omitempty"`
}

type goalsSummaryResponse struct {
	DailyCalories        *float64 `json:"daily_calories,omitempty"`
	DailyProteinGrams    *float64 `json:"daily_protein_grams,omitempty"`
	DailyCarbsGrams      *float64 `json:"daily_carbs_grams,omitempty"`
	DailyFatGrams        *float64 `json:"daily_fat_grams,omitempty"`
	DailyFiberGrams      *float64 `json:"daily_fiber_grams,omitempty"`
	DailySteps           *int     `json:"daily_steps,omitempty"`
	DailyExerciseMinutes *int     `json:"daily_exercise_minutes,omitempty"`
	TargetWeightKg       *float64 `json:"target_weight_kg,omitempty"`
}

type summaryResponse struct {
	Date     string                       `json:"date"`
	Meals    mealsSummaryResponse         `json:"meals"`
	Exercise exerciseSummaryResponse      `json:"exercise"`
	Weight   *weightEntrySummaryResponse  `json:"weight,omitempty"`
	Goals    *goalsSummaryResponse        `json:"goals,omitempty"`
}

func summaryFromDomain(s *domain.DailySummary) *summaryResponse {
	response := &summaryResponse{
		Date: s.Date.Format("2006-01-02"),
		Meals: mealsSummaryResponse{
			TotalCalories:     s.MealsSummary.TotalCalories,
			TotalProteinGrams: s.MealsSummary.TotalProteinGrams,
			TotalCarbsGrams:   s.MealsSummary.TotalCarbsGrams,
			TotalFatGrams:     s.MealsSummary.TotalFatGrams,
			TotalFiberGrams:   s.MealsSummary.TotalFiberGrams,
			Count:             s.MealsSummary.Count,
		},
		Exercise: exerciseSummaryResponse{
			TotalCaloriesBurned:  s.ExerciseSummary.TotalCaloriesBurned,
			TotalSteps:           s.ExerciseSummary.TotalSteps,
			TotalDurationSeconds: s.ExerciseSummary.TotalDurationSeconds,
			TotalDistanceMeters:  s.ExerciseSummary.TotalDistanceMeters,
			Count:                s.ExerciseSummary.Count,
		},
	}
	if s.WeightEntry != nil {
		response.Weight = &weightEntrySummaryResponse{
			WeightKg:          s.WeightEntry.WeightKg,
			BodyFatPercentage: s.WeightEntry.BodyFatPercentage,
		}
	}
	if s.Goals != nil {
		response.Goals = &goalsSummaryResponse{
			DailyCalories:        s.Goals.DailyCalories,
			DailyProteinGrams:    s.Goals.DailyProteinGrams,
			DailyCarbsGrams:      s.Goals.DailyCarbsGrams,
			DailyFatGrams:        s.Goals.DailyFatGrams,
			DailyFiberGrams:      s.Goals.DailyFiberGrams,
			DailySteps:           s.Goals.DailySteps,
			DailyExerciseMinutes: s.Goals.DailyExerciseMinutes,
			TargetWeightKg:       s.Goals.TargetWeightKg,
		}
	}
	return response
}
