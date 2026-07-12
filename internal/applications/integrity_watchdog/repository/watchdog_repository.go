package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/units"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
)

type watchdogRepository struct {
	db *gorm.DB
}

var _ ports.WatchdogRepository = (*watchdogRepository)(nil)

func NewWatchdogRepository(db *gorm.DB) *watchdogRepository {
	return &watchdogRepository{
		db: db,
	}
}

func (r *watchdogRepository) AllPhotoURLs() ([]string, error) {
	var urls []string
	err := r.db.Raw(`
		SELECT url FROM meal_photos
		UNION
		SELECT url FROM daily_photos
		UNION
		SELECT url FROM user_profile_photos
		UNION
		SELECT photo_url FROM foods WHERE photo_url != ''
	`).Scan(&urls).Error
	return urls, err
}

func (r *watchdogRepository) CrossContextPhotos() ([]ports.CrossContextPhoto, error) {
	type row struct {
		Id         uuid.UUID `gorm:"column:id"`
		MealId     uuid.UUID `gorm:"column:meal_id"`
		MealItemId uuid.UUID `gorm:"column:meal_item_id"`
	}
	var rows []row
	err := r.db.Raw(`
		SELECT mp.id, mp.meal_id, mp.meal_item_id
		FROM meal_photos mp
		JOIN meal_items mi ON mi.id = mp.meal_item_id
		WHERE mp.meal_item_id IS NOT NULL
		  AND mi.meal_id != mp.meal_id
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]ports.CrossContextPhoto, len(rows))
	for i, row := range rows {
		result[i] = ports.CrossContextPhoto{
			PhotoId:    row.Id,
			MealId:     row.MealId,
			MealItemId: row.MealItemId,
		}
	}
	return result, nil
}

func (r *watchdogRepository) MealGroupsMissingPrimary() ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.db.Raw(`
		SELECT meal_id
		FROM meal_photos
		WHERE meal_item_id IS NULL
		GROUP BY meal_id
		HAVING COUNT(*) FILTER (WHERE is_primary) = 0
	`).Scan(&ids).Error
	return ids, err
}

func (r *watchdogRepository) ItemGroupsMissingPrimary() ([]ports.ItemGroup, error) {
	type row struct {
		MealId     uuid.UUID `gorm:"column:meal_id"`
		MealItemId uuid.UUID `gorm:"column:meal_item_id"`
	}
	var rows []row
	err := r.db.Raw(`
		SELECT meal_id, meal_item_id
		FROM meal_photos
		WHERE meal_item_id IS NOT NULL
		GROUP BY meal_id, meal_item_id
		HAVING COUNT(*) FILTER (WHERE is_primary) = 0
	`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]ports.ItemGroup, len(rows))
	for i, row := range rows {
		result[i] = ports.ItemGroup{
			MealId:     row.MealId,
			MealItemId: row.MealItemId,
		}
	}
	return result, nil
}

func (r *watchdogRepository) InvalidFoodBaseUnits() ([]ports.InvalidFoodUnit, error) {
	type row struct {
		Id       uuid.UUID `gorm:"column:id"`
		BaseUnit string    `gorm:"column:base_unit"`
	}
	var rows []row
	err := r.db.Raw(`SELECT id, base_unit FROM foods`).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	var result []ports.InvalidFoodUnit
	for _, row := range rows {
		if !units.IsMetricUnit(row.BaseUnit) {
			item := ports.InvalidFoodUnit{
				FoodId:   row.Id,
				BaseUnit: row.BaseUnit,
			}
			result = append(result, item)
		}
	}
	return result, nil
}
