package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
)

type exercise struct {
	Id                      uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId                  uuid.UUID `gorm:"type:uuid;not null"`
	Date                    time.Time `gorm:"type:date;not null"`
	Type                    string    `gorm:"not null"`
	Name                    string    `gorm:"default:''"`
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
	Notes                   string `gorm:"not null;default:''"`
	ExternalId              *string
	ImportSource            *string
	CreatedAt               time.Time     `gorm:"not null;autoCreateTime"`
	UpdatedAt               time.Time     `gorm:"not null;autoUpdateTime"`
	Tags                    []exerciseTagMap `gorm:"foreignKey:ExerciseId"`
}

type exerciseTag struct {
	Id     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId uuid.UUID `gorm:"type:uuid;not null"`
	Name   string    `gorm:"not null"`
}

func (exerciseTag) TableName() string { return "exercise_tags" }

type exerciseTagMap struct {
	ExerciseId uuid.UUID   `gorm:"type:uuid;primaryKey"`
	TagId      uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Tag        exerciseTag `gorm:"foreignKey:TagId;references:Id"`
}

func (exerciseTagMap) TableName() string { return "exercise_tag_map" }

func (m *exercise) toDomain() *domain.Exercise {
	tags := make([]string, len(m.Tags))
	for i, t := range m.Tags {
		tags[i] = t.Tag.Name
	}
	return &domain.Exercise{
		Id:                      m.Id,
		UserId:                  m.UserId,
		Date:                    m.Date,
		Type:                    m.Type,
		Name:                    m.Name,
		StartedAt:               m.StartedAt,
		DurationSeconds:         m.DurationSeconds,
		EstimatedCaloriesBurned: m.EstimatedCaloriesBurned,
		Steps:                   m.Steps,
		DistanceMeters:          m.DistanceMeters,
		AverageSpeedKmh:         m.AverageSpeedKmh,
		MaxSpeedKmh:             m.MaxSpeedKmh,
		AveragePaceMinPerKm:     m.AveragePaceMinPerKm,
		ElevationGainMeters:     m.ElevationGainMeters,
		AverageHeartRate:        m.AverageHeartRate,
		MaxHeartRate:            m.MaxHeartRate,
		TotalVolumeKg:           m.TotalVolumeKg,
		TotalSets:               m.TotalSets,
		Tags:                    tags,
		Notes:                   m.Notes,
		ExternalId:              m.ExternalId,
		ImportSource:            m.ImportSource,
		CreatedAt:               m.CreatedAt,
		UpdatedAt:               m.UpdatedAt,
	}
}

func exerciseFromDomain(e *domain.Exercise) *exercise {
	return &exercise{
		Id:                      e.Id,
		UserId:                  e.UserId,
		Date:                    e.Date,
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
		Notes:                   e.Notes,
		ExternalId:              e.ExternalId,
		ImportSource:            e.ImportSource,
		CreatedAt:               e.CreatedAt,
		UpdatedAt:               e.UpdatedAt,
	}
}
