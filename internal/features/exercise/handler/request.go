package handler

import "time"

type createExerciseRequest struct {
	Date                    string     `json:"date"`
	Type                    string     `json:"type"`
	Name                    string     `json:"name"`
	StartedAt               *time.Time `json:"started_at,omitempty"`
	DurationSeconds         *int       `json:"duration_seconds,omitempty"`
	EstimatedCaloriesBurned *float64   `json:"estimated_calories_burned,omitempty"`
	Steps                   *int       `json:"steps,omitempty"`
	DistanceMeters          *float64   `json:"distance_meters,omitempty"`
	MaxSpeedKmh             *float64   `json:"max_speed_kmh,omitempty"`
	ElevationGainMeters     *float64   `json:"elevation_gain_meters,omitempty"`
	AverageHeartRate        *int       `json:"average_heart_rate,omitempty"`
	MaxHeartRate            *int       `json:"max_heart_rate,omitempty"`
	TotalVolumeKg           *float64   `json:"total_volume_kg,omitempty"`
	TotalSets               *int       `json:"total_sets,omitempty"`
	Tags                    []string   `json:"tags"`
	Notes                   string     `json:"notes"`
}

type updateExerciseRequest struct {
	Date                    *string    `json:"date,omitempty"`
	Type                    *string    `json:"type,omitempty"`
	Name                    *string    `json:"name,omitempty"`
	StartedAt               *time.Time `json:"started_at,omitempty"`
	DurationSeconds         *int       `json:"duration_seconds,omitempty"`
	EstimatedCaloriesBurned *float64   `json:"estimated_calories_burned,omitempty"`
	Steps                   *int       `json:"steps,omitempty"`
	DistanceMeters          *float64   `json:"distance_meters,omitempty"`
	MaxSpeedKmh             *float64   `json:"max_speed_kmh,omitempty"`
	ElevationGainMeters     *float64   `json:"elevation_gain_meters,omitempty"`
	AverageHeartRate        *int       `json:"average_heart_rate,omitempty"`
	MaxHeartRate            *int       `json:"max_heart_rate,omitempty"`
	TotalVolumeKg           *float64   `json:"total_volume_kg,omitempty"`
	TotalSets               *int       `json:"total_sets,omitempty"`
	Tags                    *[]string  `json:"tags,omitempty"`
	Notes                   *string    `json:"notes,omitempty"`
}
