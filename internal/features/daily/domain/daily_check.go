package domain

import "time"

type DailyCheck struct {
	Date                time.Time
	Complete            bool
	MissingMeasurements int  // meal items without measurement_method
	MealsWithoutPhoto   int  // meals that have no photo at all
	HasDailyPhoto       bool
	HasSteps            bool
	HasExercise         bool
	HasRecentWeight     bool // weight entry within the last 7 days
}
