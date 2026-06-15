package domain

import "errors"

type ActivityType string

const ActivityTypeSteps ActivityType = "steps"

var (
	ErrUnsupportedActivityType = errors.New("unsupported activity type")
	ErrNoWeightData            = errors.New("no weight data available; log a weight entry first")
)

type EstimateRequest struct {
	Type  ActivityType
	Value float64
}

type EstimateResult struct {
	Type              ActivityType
	Value             float64
	EstimatedCalories float64
	WeightKg          float64
}
