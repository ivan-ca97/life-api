package service

import (
	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/dayclosure"
	"github.com/ivan-ca97/life/pkg/types"
	"github.com/ivan-ca97/life/pkg/validate"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
	"github.com/ivan-ca97/life/internal/features/exercise/ports"
)

type exerciseService struct {
	repository     ports.ExerciseRepository
	closureChecker dayclosure.DayClosureChecker
}

var _ ports.ExerciseService = (*exerciseService)(nil)

func NewExerciseService(repository ports.ExerciseRepository, closureChecker dayclosure.DayClosureChecker) *exerciseService {
	return &exerciseService{
		repository:     repository,
		closureChecker: closureChecker,
	}
}

func (s *exerciseService) Create(userId uuid.UUID, params ports.CreateParams) (*domain.Exercise, error) {
	closed, err := s.closureChecker.IsClosed(userId, params.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	err = validate.NonEmpty(params.Name, "name")
	if err != nil {
		return nil, err
	}
	if !domain.IsValidExerciseType(params.Type) {
		return nil, domain.ErrInvalidExerciseType
	}
	err = validate.NonNegativeIntPtr(params.DurationSeconds, "duration_seconds")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.Steps, "steps")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.DistanceMeters, "distance_meters")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.EstimatedCaloriesBurned, "estimated_calories_burned")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.ElevationGainMeters, "elevation_gain_meters")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.AverageHeartRate, "average_heart_rate")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.MaxHeartRate, "max_heart_rate")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.TotalVolumeKg, "total_volume_kg")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.TotalSets, "total_sets")
	if err != nil {
		return nil, err
	}
	avgSpeed, avgPace := computeSpeedAndPace(params.DistanceMeters, params.DurationSeconds)
	exercise := &domain.Exercise{
		Id:                      uuid.New(),
		UserId:                  userId,
		Date:                    params.Date,
		Type:                    params.Type,
		Name:                    params.Name,
		StartedAt:               params.StartedAt,
		DurationSeconds:         params.DurationSeconds,
		EstimatedCaloriesBurned: params.EstimatedCaloriesBurned,
		Steps:                   params.Steps,
		DistanceMeters:          params.DistanceMeters,
		AverageSpeedKmh:         avgSpeed,
		MaxSpeedKmh:             params.MaxSpeedKmh,
		AveragePaceMinPerKm:     avgPace,
		ElevationGainMeters:     params.ElevationGainMeters,
		AverageHeartRate:        params.AverageHeartRate,
		MaxHeartRate:            params.MaxHeartRate,
		TotalVolumeKg:           params.TotalVolumeKg,
		TotalSets:               params.TotalSets,
		Tags:                    params.Tags,
		Notes:                   params.Notes,
		ExternalId:              params.ExternalId,
		ImportSource:            params.ImportSource,
	}
	err = s.repository.Create(exercise)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *exerciseService) GetById(id, userId uuid.UUID) (*domain.Exercise, error) {
	exercise, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *exerciseService) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.Exercise], error) {
	page, err := s.repository.List(userId, params)
	if err != nil {
		return types.Page[domain.Exercise]{}, err
	}
	return page, nil
}

func (s *exerciseService) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.Exercise, error) {
	current, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	closed, err := s.closureChecker.IsClosed(userId, current.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	if params.Type != nil && !domain.IsValidExerciseType(*params.Type) {
		return nil, domain.ErrInvalidExerciseType
	}
	err = validate.NonNegativeIntPtr(params.DurationSeconds, "duration_seconds")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.Steps, "steps")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.DistanceMeters, "distance_meters")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.EstimatedCaloriesBurned, "estimated_calories_burned")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.ElevationGainMeters, "elevation_gain_meters")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.AverageHeartRate, "average_heart_rate")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.MaxHeartRate, "max_heart_rate")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativePtr(params.TotalVolumeKg, "total_volume_kg")
	if err != nil {
		return nil, err
	}
	err = validate.NonNegativeIntPtr(params.TotalSets, "total_sets")
	if err != nil {
		return nil, err
	}
	if params.DistanceMeters != nil || params.DurationSeconds != nil {
		distance := current.DistanceMeters
		if params.DistanceMeters != nil {
			distance = params.DistanceMeters
		}
		duration := current.DurationSeconds
		if params.DurationSeconds != nil {
			duration = params.DurationSeconds
		}
		avgSpeed, avgPace := computeSpeedAndPace(distance, duration)
		params.AverageSpeedKmh = avgSpeed
		params.AveragePaceMinPerKm = avgPace
	}
	exercise, err := s.repository.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func computeSpeedAndPace(distance *float64, duration *int) (*float64, *float64) {
	if distance == nil || duration == nil || *distance <= 0 || *duration <= 0 {
		return nil, nil
	}
	d := *distance
	s := float64(*duration)
	avgSpeed := (d / 1000) / (s / 3600)
	avgPace := (s / 60) / (d / 1000)
	return &avgSpeed, &avgPace
}

func (s *exerciseService) Delete(id, userId uuid.UUID) error {
	exercise, err := s.repository.FindById(id, userId)
	if err != nil {
		return err
	}
	closed, err := s.closureChecker.IsClosed(userId, exercise.Date)
	if err != nil {
		return err
	}
	if closed {
		return dayclosure.ErrDayClosed
	}

	err = s.repository.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}
