package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	ExerciseTypeWeightlifting    = "weightlifting"
	ExerciseTypeWalking          = "walking"
	ExerciseTypeCycling          = "cycling"
	ExerciseTypeRunning          = "running"
	ExerciseTypeOther            = "other"
	ExerciseTypeManualAdjustment = "manual_adjustment"
)

const (
	ImportSourceHealthConnect     = "health_connect"
	ImportSourceHevy              = "hevy"
	ImportSourceHealthConnectHevy = "health_connect+hevy"
)

func IsValidExerciseType(exerciseType string) bool {
	switch exerciseType {
	case ExerciseTypeWeightlifting, ExerciseTypeWalking, ExerciseTypeCycling,
		ExerciseTypeRunning, ExerciseTypeOther, ExerciseTypeManualAdjustment:
		return true
	}
	return false
}

type Exercise struct {
	Id                      uuid.UUID
	UserId                  uuid.UUID
	Date                    time.Time
	Type                    string
	Name                    string
	StartedAt               *time.Time
	DurationSeconds         *int
	EstimatedCaloriesBurned *float64
	Steps                   *int
	DistanceMeters          *float64
	AverageSpeedKmh         *float64
	MaxSpeedKmh             *float64
	AveragePaceMinPerKm     *float64
	ElevationGainMeters     *float64
	AverageHeartRate        *int
	MaxHeartRate            *int
	TotalVolumeKg           *float64
	TotalSets               *int
	Tags                    []string
	Notes                   string
	ExternalId              *string
	ImportSource            *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
