package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
)

type exerciseResponse struct {
	Id                      uuid.UUID  `json:"id"`
	Date                    string     `json:"date"`
	Type                    string     `json:"type"`
	Name                    string     `json:"name"`
	StartedAt               *time.Time `json:"started_at,omitempty"`
	DurationSeconds         *int       `json:"duration_seconds,omitempty"`
	EstimatedCaloriesBurned *float64   `json:"estimated_calories_burned,omitempty"`
	Steps                   *int       `json:"steps,omitempty"`
	DistanceMeters          *float64   `json:"distance_meters,omitempty"`
	AverageSpeedKmh         *float64   `json:"average_speed_kmh,omitempty"`
	MaxSpeedKmh             *float64   `json:"max_speed_kmh,omitempty"`
	AveragePaceMinPerKm     *float64   `json:"average_pace_min_per_km,omitempty"`
	ElevationGainMeters     *float64   `json:"elevation_gain_meters,omitempty"`
	AverageHeartRate        *int       `json:"average_heart_rate,omitempty"`
	MaxHeartRate            *int       `json:"max_heart_rate,omitempty"`
	TotalVolumeKg           *float64   `json:"total_volume_kg,omitempty"`
	TotalSets               *int       `json:"total_sets,omitempty"`
	Tags                    []string   `json:"tags"`
	Notes                   string     `json:"notes"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

func exerciseFromDomain(e *domain.Exercise) *exerciseResponse {
	tags := e.Tags
	if tags == nil {
		tags = []string{}
	}
	return &exerciseResponse{
		Id:                      e.Id,
		Date:                    e.Date.Format("2006-01-02"),
		Type:                    e.Type,
		Name:                    e.Name,
		StartedAt:               e.StartedAt,
		DurationSeconds:         e.DurationSeconds,
		EstimatedCaloriesBurned: e.EstimatedCaloriesBurned,
		Steps:                   e.Steps,
		DistanceMeters:          e.DistanceMeters,
		AverageSpeedKmh:         e.AverageSpeedKmh,
		MaxSpeedKmh:             e.MaxSpeedKmh,
		AveragePaceMinPerKm:     e.AveragePaceMinPerKm,
		ElevationGainMeters:     e.ElevationGainMeters,
		AverageHeartRate:        e.AverageHeartRate,
		MaxHeartRate:            e.MaxHeartRate,
		TotalVolumeKg:           e.TotalVolumeKg,
		TotalSets:               e.TotalSets,
		Tags:                    tags,
		Notes:                   e.Notes,
		CreatedAt:               e.CreatedAt,
		UpdatedAt:               e.UpdatedAt,
	}
}

type exercisePage struct {
	Items  []exerciseResponse `json:"items"`
	Total  int64              `json:"total"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

func newExercisePage(page types.Page[domain.Exercise]) *exercisePage {
	items := make([]exerciseResponse, len(page.Items))
	for i, e := range page.Items {
		items[i] = *exerciseFromDomain(&e)
	}
	return &exercisePage{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}
