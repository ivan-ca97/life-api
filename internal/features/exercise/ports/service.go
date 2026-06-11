package ports

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
)

type CreateParams struct {
	Date                    time.Time
	Type                    string
	Name                    string
	StartedAt               *time.Time
	DurationSeconds         *int
	EstimatedCaloriesBurned *float64
	Steps                   *int
	DistanceMeters          *float64
	MaxSpeedKmh             *float64
	ElevationGainMeters     *float64
	AverageHeartRate        *int
	MaxHeartRate            *int
	TotalVolumeKg           *float64
	TotalSets               *int
	Tags                    []string
	Notes                   string
	ExternalId              *string
}

type UpdateParams struct {
	Date                    *time.Time
	Type                    *string
	Name                    *string
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
	Tags                    *[]string
	Notes                   *string
}

type ListParams struct {
	types.PaginationParams
	Date *time.Time
}

type ExerciseService interface {
	Create(userId uuid.UUID, params CreateParams) (*domain.Exercise, error)
	GetById(id, userId uuid.UUID) (*domain.Exercise, error)
	List(userId uuid.UUID, params ListParams) (types.Page[domain.Exercise], error)
	Update(id, userId uuid.UUID, params UpdateParams) (*domain.Exercise, error)
	Delete(id, userId uuid.UUID) error
}
