package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
	"github.com/ivan-ca97/life/internal/features/meal/ports"
)

type mealRepository struct {
	db *gorm.DB
}

var _ ports.MealRepository = (*mealRepository)(nil)

func NewMealRepository(db *gorm.DB) *mealRepository {
	return &mealRepository{
		db: db,
	}
}

func (r *mealRepository) Create(m *domain.Meal) error {
	model := mealFromDomain(m)
	err := r.db.Omit("Tags", "Items").Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting meal", err)
	}
	if len(m.Tags) > 0 {
		tags := make([]mealTag, len(m.Tags))
		for i, t := range m.Tags {
			tags[i] = mealTag{
				MealId: m.Id,
				Tag:    t,
			}
		}
		err = r.db.Create(&tags).Error
		if err != nil {
			return cerr.NewInternalError("inserting meal tags", err)
		}
	}
	if len(m.Items) > 0 {
		items := make([]mealItem, len(m.Items))
		for i, item := range m.Items {
			var inputUnit *string
			if item.InputUnit != "" {
				inputUnit = &item.InputUnit
			}
			var normalizedUnit *string
			if item.NormalizedUnit != "" {
				normalizedUnit = &item.NormalizedUnit
			}
			var normalizedQty *float64
			if item.NormalizedQuantity != 0 {
				nq := item.NormalizedQuantity
				normalizedQty = &nq
			}
			var method *string
			if item.MeasurementMethod != "" {
				m := string(item.MeasurementMethod)
				method = &m
			}
			items[i] = mealItem{
				Id:                 uuid.New(),
				MealId:             m.Id,
				FoodId:             item.FoodId,
				InputQuantity:      item.InputQuantity,
				InputUnit:          inputUnit,
				NormalizedQuantity: normalizedQty,
				NormalizedUnit:     normalizedUnit,
				Calories:           item.Calories,
				ProteinGrams:       item.ProteinGrams,
				CarbsGrams:         item.CarbsGrams,
				FatGrams:           item.FatGrams,
				FiberGrams:         item.FiberGrams,
				Notes:              item.Notes,
				MeasurementMethod:  method,
			}
		}
		err = r.db.Create(&items).Error
		if err != nil {
			return cerr.NewInternalError("inserting meal items", err)
		}
	}
	return nil
}

func (r *mealRepository) FindById(id, userId uuid.UUID) (*domain.Meal, error) {
	var model meal
	err := r.db.
		Preload("Tags").
		Preload("Items").
		Preload("Items.Food").
		Where("id = ? AND user_id = ?", id, userId).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrMealNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding meal by id", err)
	}
	return model.toDomain(), nil
}

func (r *mealRepository) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.Meal], error) {
	var models []meal
	var total int64

	countQuery := r.db.Model(&meal{}).Where("user_id = ?", userId)
	if params.Date != nil {
		countQuery = countQuery.Where("date = ?", *params.Date)
	}
	err := countQuery.Count(&total).Error
	if err != nil {
		return types.Page[domain.Meal]{}, cerr.NewInternalError("counting meals", err)
	}

	findQuery := r.db.Preload("Tags").Preload("Items").Preload("Items.Food").Where("user_id = ?", userId)
	if params.Date != nil {
		findQuery = findQuery.Where("date = ?", *params.Date)
	}
	err = findQuery.
		Limit(params.Limit).
		Offset(params.Offset).
		Order("date DESC, eaten_at DESC NULLS LAST").
		Find(&models).
		Error
	if err != nil {
		return types.Page[domain.Meal]{}, cerr.NewInternalError("listing meals", err)
	}

	meals := make([]domain.Meal, len(models))
	for i, m := range models {
		meals[i] = *m.toDomain()
	}

	result := types.Page[domain.Meal]{
		Items:  meals,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	return result, nil
}

func (r *mealRepository) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.Meal, error) {
	var count int64
	err := r.db.Model(&meal{}).Where("id = ? AND user_id = ?", id, userId).Count(&count).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking meal existence", err)
	}
	if count == 0 {
		return nil, domain.ErrMealNotFound
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
	if params.PhotoUrl != nil {
		updates["photo_url"] = *params.PhotoUrl
	}
	if params.EatenAt != nil {
		updates["eaten_at"] = *params.EatenAt
	}
	if params.Calories != nil {
		updates["calories"] = *params.Calories
	}
	if params.ProteinGrams != nil {
		updates["protein_grams"] = *params.ProteinGrams
	}
	if params.CarbsGrams != nil {
		updates["carbs_grams"] = *params.CarbsGrams
	}
	if params.FatGrams != nil {
		updates["fat_grams"] = *params.FatGrams
	}
	if params.FiberGrams != nil {
		updates["fiber_grams"] = *params.FiberGrams
	}
	if params.Notes != nil {
		updates["notes"] = *params.Notes
	}

	if len(updates) > 0 {
		err = r.db.Model(&meal{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
		if err != nil {
			return nil, cerr.NewInternalError("updating meal", err)
		}
	}

	if params.Tags != nil {
		err = r.db.Where("meal_id = ?", id).Delete(&mealTag{}).Error
		if err != nil {
			return nil, cerr.NewInternalError("deleting meal tags", err)
		}
		if len(*params.Tags) > 0 {
			tags := make([]mealTag, len(*params.Tags))
			for i, t := range *params.Tags {
				tags[i] = mealTag{
					MealId: id,
					Tag:    t,
				}
			}
			err = r.db.Create(&tags).Error
			if err != nil {
				return nil, cerr.NewInternalError("inserting meal tags", err)
			}
		}
	}

	if params.ResolvedItems != nil {
		err = r.db.Where("meal_id = ?", id).Delete(&mealItem{}).Error
		if err != nil {
			return nil, cerr.NewInternalError("deleting meal items", err)
		}
		if len(*params.ResolvedItems) > 0 {
			items := make([]mealItem, len(*params.ResolvedItems))
			for i, item := range *params.ResolvedItems {
				var inputUnit *string
				if item.InputUnit != "" {
					inputUnit = &item.InputUnit
				}
				var normalizedUnit *string
				if item.NormalizedUnit != "" {
					normalizedUnit = &item.NormalizedUnit
				}
				var normalizedQty *float64
				if item.NormalizedQuantity != 0 {
					nq := item.NormalizedQuantity
					normalizedQty = &nq
				}
				var method *string
				if item.MeasurementMethod != "" {
					mv := string(item.MeasurementMethod)
					method = &mv
				}
				items[i] = mealItem{
					Id:                 uuid.New(),
					MealId:             id,
					FoodId:             item.FoodId,
					InputQuantity:      item.InputQuantity,
					InputUnit:          inputUnit,
					NormalizedQuantity: normalizedQty,
					NormalizedUnit:     normalizedUnit,
					Calories:           item.Calories,
					ProteinGrams:       item.ProteinGrams,
					CarbsGrams:         item.CarbsGrams,
					FatGrams:           item.FatGrams,
					FiberGrams:         item.FiberGrams,
					Notes:              item.Notes,
					MeasurementMethod:  method,
				}
			}
			err = r.db.Create(&items).Error
			if err != nil {
				return nil, cerr.NewInternalError("inserting meal items", err)
			}
		}
	}

	return r.FindById(id, userId)
}

func (r *mealRepository) ListDistinctTypes(userId uuid.UUID, hour *int) ([]string, error) {
	type typeRow struct {
		Type string
	}
	var rows []typeRow

	if hour != nil {
		query := `
			SELECT type
			FROM meals
			WHERE user_id = ?
			GROUP BY type
			ORDER BY
				COUNT(*) FILTER (WHERE eaten_at IS NOT NULL AND LEAST(
					ABS(EXTRACT(HOUR FROM eaten_at) - ?),
					24 - ABS(EXTRACT(HOUR FROM eaten_at) - ?)
				) <= 2) DESC,
				COUNT(*) DESC
		`
		err := r.db.Raw(query, userId, *hour, *hour).Scan(&rows).Error
		if err != nil {
			return nil, cerr.NewInternalError("listing meal types", err)
		}
	} else {
		query := `
			SELECT type
			FROM meals
			WHERE user_id = ?
			GROUP BY type
			ORDER BY COUNT(*) DESC
		`
		err := r.db.Raw(query, userId).Scan(&rows).Error
		if err != nil {
			return nil, cerr.NewInternalError("listing meal types", err)
		}
	}

	types := make([]string, len(rows))
	for i, row := range rows {
		types[i] = row.Type
	}
	return types, nil
}

func (r *mealRepository) Delete(id, userId uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&meal{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting meal", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrMealNotFound
	}
	return nil
}
