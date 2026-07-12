package domain

import "time"

type GoalMetric struct {
	Target      float64
	Average     float64
	DaysMet     int
	DaysTracked int
	DaysTotal   int
}

type WeightProgress struct {
	TargetKg  float64
	CurrentKg *float64
}

type GoalProgress struct {
	From                 time.Time
	To                   time.Time
	Goal                 *Goal
	DaysTotal            int
	DailyCalories        *GoalMetric
	DailyProteinGrams    *GoalMetric
	DailyCarbsGrams      *GoalMetric
	DailyFatGrams        *GoalMetric
	DailyFiberGrams      *GoalMetric
	DailySteps           *GoalMetric
	DailyExerciseMinutes *GoalMetric
	WeightProgress       *WeightProgress
}
