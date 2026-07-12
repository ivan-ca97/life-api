package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/exercise/domain"
	"github.com/ivan-ca97/life/internal/features/exercise/ports"
)

type exerciseRepository struct {
	db *gorm.DB
}

var _ ports.ExerciseRepository = (*exerciseRepository)(nil)

func NewExerciseRepository(db *gorm.DB) *exerciseRepository {
	return &exerciseRepository{
		db: db,
	}
}

func (r *exerciseRepository) Create(e *domain.Exercise) error {
	model := exerciseFromDomain(e)
	err := r.db.Omit("Tags").Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting exercise", err)
	}
	if len(e.Tags) > 0 {
		err := r.upsertTags(e.Id, e.UserId, e.Tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *exerciseRepository) FindById(id, userId uuid.UUID) (*domain.Exercise, error) {
	var model exercise
	err := r.db.
		Preload("Tags.Tag").
		Where("id = ? AND user_id = ?", id, userId).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrExerciseNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding exercise by id", err)
	}
	return model.toDomain(), nil
}

func (r *exerciseRepository) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.Exercise], error) {
	var models []exercise
	var total int64

	countQuery := r.db.Model(&exercise{}).Where("user_id = ?", userId)
	if params.Date != nil {
		countQuery = countQuery.Where("date = ?", *params.Date)
	}
	err := countQuery.Count(&total).Error
	if err != nil {
		return types.Page[domain.Exercise]{}, cerr.NewInternalError("counting exercises", err)
	}

	findQuery := r.db.Preload("Tags.Tag").Where("user_id = ?", userId)
	if params.Date != nil {
		findQuery = findQuery.Where("date = ?", *params.Date)
	}
	err = findQuery.
		Limit(params.Limit).
		Offset(params.Offset).
		Order("date DESC, started_at DESC NULLS LAST").
		Find(&models).
		Error
	if err != nil {
		return types.Page[domain.Exercise]{}, cerr.NewInternalError("listing exercises", err)
	}

	exercises := make([]domain.Exercise, len(models))
	for i, m := range models {
		exercises[i] = *m.toDomain()
	}

	result := types.Page[domain.Exercise]{
		Items:  exercises,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	return result, nil
}

func (r *exerciseRepository) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.Exercise, error) {
	var count int64
	err := r.db.Model(&exercise{}).Where("id = ? AND user_id = ?", id, userId).Count(&count).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking exercise existence", err)
	}
	if count == 0 {
		return nil, domain.ErrExerciseNotFound
	}

	updates := map[string]any{}
	if params.Date != nil {
		updates["date"] = *params.Date
	}
	if params.Type != nil {
		updates["type"] = *params.Type
	}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.StartedAt != nil {
		updates["started_at"] = *params.StartedAt
	}
	if params.DurationSeconds != nil {
		updates["duration_seconds"] = *params.DurationSeconds
	}
	if params.EstimatedCaloriesBurned != nil {
		updates["estimated_calories_burned"] = *params.EstimatedCaloriesBurned
	}
	if params.Steps != nil {
		updates["steps"] = *params.Steps
	}
	if params.DistanceMeters != nil {
		updates["distance_meters"] = *params.DistanceMeters
	}
	if params.AverageSpeedKmh != nil {
		updates["average_speed_kmh"] = *params.AverageSpeedKmh
	}
	if params.MaxSpeedKmh != nil {
		updates["max_speed_kmh"] = *params.MaxSpeedKmh
	}
	if params.AveragePaceMinPerKm != nil {
		updates["average_pace_min_per_km"] = *params.AveragePaceMinPerKm
	}
	if params.ElevationGainMeters != nil {
		updates["elevation_gain_meters"] = *params.ElevationGainMeters
	}
	if params.AverageHeartRate != nil {
		updates["average_heart_rate"] = *params.AverageHeartRate
	}
	if params.MaxHeartRate != nil {
		updates["max_heart_rate"] = *params.MaxHeartRate
	}
	if params.TotalVolumeKg != nil {
		updates["total_volume_kg"] = *params.TotalVolumeKg
	}
	if params.TotalSets != nil {
		updates["total_sets"] = *params.TotalSets
	}
	if params.Notes != nil {
		updates["notes"] = *params.Notes
	}
	if params.ImportSource != nil {
		updates["import_source"] = *params.ImportSource
	}

	if len(updates) > 0 {
		err = r.db.Model(&exercise{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
		if err != nil {
			return nil, cerr.NewInternalError("updating exercise", err)
		}
	}

	if params.Tags != nil {
		if err = r.db.Where("exercise_id = ?", id).Delete(&exerciseTagMap{}).Error; err != nil {
			return nil, cerr.NewInternalError("deleting exercise tag map", err)
		}
		if len(*params.Tags) > 0 {
			if err = r.upsertTags(id, userId, *params.Tags); err != nil {
				return nil, err
			}
		}
	}

	return r.FindById(id, userId)
}

func (r *exerciseRepository) upsertTags(exerciseId, userId uuid.UUID, names []string) error {
	entries := make([]exerciseTag, len(names))
	for i, name := range names {
		entries[i] = exerciseTag{
			Id:     uuid.New(),
			UserId: userId,
			Name:   name,
		}
	}
	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "name"}},
		DoNothing: true,
	}).Create(&entries).Error; err != nil {
		return cerr.NewInternalError("upserting exercise tags", err)
	}
	var tags []exerciseTag
	err := r.db.Where("user_id = ? AND name IN ?", userId, names).Find(&tags).Error
	if err != nil {
		return cerr.NewInternalError("fetching exercise tag ids", err)
	}
	maps := make([]exerciseTagMap, len(tags))
	for i, t := range tags {
		maps[i] = exerciseTagMap{
			ExerciseId: exerciseId,
			TagId:      t.Id,
		}
	}
	return r.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&maps).Error
}

func (r *exerciseRepository) Delete(id, userId uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&exercise{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting exercise", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrExerciseNotFound
	}
	return nil
}

func (r *exerciseRepository) ExistsByDateAndName(userId uuid.UUID, date time.Time, name string) (bool, error) {
	var count int64
	err := r.db.
		Model(&exercise{}).
		Where("user_id = ? AND date = ? AND name = ?", userId, date, name).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking exercise existence", err)
	}
	return count > 0, nil
}

func (r *exerciseRepository) FindByDateAndName(userId uuid.UUID, date time.Time, name string) (*domain.Exercise, error) {
	var model exercise
	err := r.db.
		Preload("Tags.Tag").
		Where("user_id = ? AND date = ? AND name = ?", userId, date, name).
		First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrExerciseNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding exercise by date and name", err)
	}
	return model.toDomain(), nil
}

func (r *exerciseRepository) ExistsByExternalId(userId uuid.UUID, externalId string) (bool, error) {
	var count int64
	err := r.db.
		Model(&exercise{}).
		Where("user_id = ? AND external_id = ?", userId, externalId).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking exercise existence by external id", err)
	}
	return count > 0, nil
}
