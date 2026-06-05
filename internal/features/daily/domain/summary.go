package domain

import "time"

type DailySummary struct {
	Date            time.Time
	MealsSummary    MealsSummary
	ExerciseSummary ExerciseSummary
	WeightEntry     *WeightEntrySummary
	Goals           *GoalsSummary
	EstimatedBMR    *float64
	CaloricBalance  *float64
}

type MealsSummary struct {
	TotalCalories     float64
	TotalProteinGrams float64
	TotalCarbsGrams   float64
	TotalFatGrams     float64
	TotalFiberGrams   float64
	Count             int
}

type ExerciseSummary struct {
	TotalCaloriesBurned  float64
	TotalSteps           int
	TotalDurationSeconds int
	TotalDistanceMeters  float64
	Count                int
}

type WeightEntrySummary struct {
	WeightKg          float64
	BodyFatPercentage *float64
}

type GoalsSummary struct {
	DailyCalories        *float64
	DailyProteinGrams    *float64
	DailyCarbsGrams      *float64
	DailyFatGrams        *float64
	DailyFiberGrams      *float64
	DailySteps           *int
	DailyExerciseMinutes *int
	TargetWeightKg       *float64
}
