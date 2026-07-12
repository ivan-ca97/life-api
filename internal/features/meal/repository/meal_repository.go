package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	err := r.db.Omit("Tags", "Items", "Photos").Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting meal", err)
	}

	// Items are created before photos so that ItemFoodId references can be resolved.
	foodIdToItemId := make(map[uuid.UUID]uuid.UUID, len(m.Items))
	if len(m.Items) > 0 {
		items := make([]mealItem, len(m.Items))
		for i, item := range m.Items {
			mi := buildMealItem(m.Id, item)
			mi.Id = uuid.New()
			items[i] = mi
			foodIdToItemId[item.FoodId] = mi.Id
		}
		err = r.db.Create(&items).Error
		if err != nil {
			return cerr.NewInternalError("inserting meal items", err)
		}
	}

	if len(m.Photos) > 0 {
		photos := make([]mealPhoto, len(m.Photos))
		for i, p := range m.Photos {
			photos[i] = mealPhoto{
				Id:         p.Id,
				MealId:     m.Id,
				MealItemId: resolveItemId(p.MealItemId, p.ItemFoodId, foodIdToItemId),
				Url:        p.Url,
				IsPrimary:  p.IsPrimary,
			}
		}
		err = r.db.Create(&photos).Error
		if err != nil {
			return cerr.NewInternalError("inserting meal photos", err)
		}
	}

	if len(m.Tags) > 0 {
		err := r.upsertTags(m.Id, m.UserId, m.Tags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *mealRepository) FindById(id, userId uuid.UUID) (*domain.Meal, error) {
	var model meal
	err := r.db.
		Preload("Tags.Tag").
		Preload("Items").
		Preload("Items.Food").
		Preload("Photos").
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
	countQuery = applyMealFilters(countQuery, params)
	err := countQuery.Count(&total).Error
	if err != nil {
		return types.Page[domain.Meal]{}, cerr.NewInternalError("counting meals", err)
	}

	findQuery := r.db.Preload("Tags.Tag").Preload("Items").Preload("Items.Food").Preload("Photos").Where("user_id = ?", userId)
	findQuery = applyMealFilters(findQuery, params)
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
	if params.Status != nil {
		updates["status"] = *params.Status
	}

	if len(updates) > 0 {
		err = r.db.Model(&meal{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
		if err != nil {
			return nil, cerr.NewInternalError("updating meal", err)
		}
	}

	if params.Tags != nil {
		if err = r.db.Where("meal_id = ?", id).Delete(&mealTagMap{}).Error; err != nil {
			return nil, cerr.NewInternalError("deleting meal tag map", err)
		}
		if len(*params.Tags) > 0 {
			if err = r.upsertTags(id, userId, *params.Tags); err != nil {
				return nil, err
			}
		}
	}

	// Items are upserted before photos so that photo ItemFoodId references can be resolved.
	var foodIdToItemId map[uuid.UUID]uuid.UUID
	if params.ResolvedItems != nil {
		foodIdToItemId, err = r.upsertItems(id, *params.ResolvedItems)
		if err != nil {
			return nil, err
		}
	}

	if params.Photos != nil {
		err = r.db.Where("meal_id = ?", id).Delete(&mealPhoto{}).Error
		if err != nil {
			return nil, cerr.NewInternalError("deleting meal photos", err)
		}
		if len(*params.Photos) > 0 {
			photos := make([]mealPhoto, len(*params.Photos))
			for i, p := range *params.Photos {
				photos[i] = mealPhoto{
					Id:         uuid.New(),
					MealId:     id,
					MealItemId: resolveItemId(p.MealItemId, p.ItemFoodId, foodIdToItemId),
					Url:        p.Url,
					IsPrimary:  p.IsPrimary,
				}
			}
			err = r.db.Create(&photos).Error
			if err != nil {
				return nil, cerr.NewInternalError("inserting meal photos", err)
			}
		}
	}

	updated, err := r.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// upsertItems preserves item UUIDs for food_ids already in the meal (update),
// assigns new UUIDs for food_ids being added, and deletes food_ids no longer present.
// Returns a food_id → item.Id map for resolving photo ItemFoodId references.
func (r *mealRepository) upsertItems(mealId uuid.UUID, incoming []domain.MealItem) (map[uuid.UUID]uuid.UUID, error) {
	var existingItems []mealItem
	err := r.db.Where("meal_id = ?", mealId).Find(&existingItems).Error
	if err != nil {
		return nil, cerr.NewInternalError("fetching meal items for upsert", err)
	}

	existingByFoodId := make(map[uuid.UUID]uuid.UUID, len(existingItems))
	for _, item := range existingItems {
		existingByFoodId[item.FoodId] = item.Id
	}

	incomingFoodIds := make(map[uuid.UUID]bool, len(incoming))
	newItems := make([]mealItem, 0, len(incoming))
	foodIdToItemId := make(map[uuid.UUID]uuid.UUID, len(incoming))

	for _, item := range incoming {
		incomingFoodIds[item.FoodId] = true
		mi := buildMealItem(mealId, item)
		existingId, ok := existingByFoodId[item.FoodId]
		if ok {
			mi.Id = existingId
		} else {
			mi.Id = uuid.New()
		}
		newItems = append(newItems, mi)
		foodIdToItemId[item.FoodId] = mi.Id
	}

	for _, item := range existingItems {
		if !incomingFoodIds[item.FoodId] {
			err := r.db.Delete(&mealItem{}, "id = ?", item.Id).Error
			if err != nil {
				return nil, cerr.NewInternalError("deleting removed meal item", err)
			}
		}
	}

	if len(newItems) > 0 {
		onConflict := clause.OnConflict{
			UpdateAll: true,
		}
		err := r.db.Clauses(onConflict).Create(&newItems).Error
		if err != nil {
			return nil, cerr.NewInternalError("upserting meal items", err)
		}
	}

	return foodIdToItemId, nil
}

// resolveItemId returns the final meal_item_id for a photo.
// MealItemId takes precedence (existing item UUID from a GET response).
// If only ItemFoodId is set, it is resolved using the food_id → item.Id map.
func resolveItemId(mealItemId, itemFoodId *uuid.UUID, foodIdToItemId map[uuid.UUID]uuid.UUID) *uuid.UUID {
	if mealItemId != nil {
		return mealItemId
	}
	if itemFoodId != nil {
		resolved, ok := foodIdToItemId[*itemFoodId]
		if ok {
			return &resolved
		}
	}
	return nil
}

func (r *mealRepository) ListDistinctTypes(userId uuid.UUID, hour *int) ([]string, error) {
	type typeRow struct {
		Type string
	}
	var rows []typeRow

	if hour != nil {
		query := `
			SELECT INITCAP(LOWER(MIN(type))) AS type
			FROM meals
			WHERE user_id = ?
			GROUP BY LOWER(type)
			ORDER BY
				COUNT(*) FILTER (WHERE LEAST(
					ABS(EXTRACT(HOUR FROM COALESCE(eaten_at, created_at)) - ?),
					24 - ABS(EXTRACT(HOUR FROM COALESCE(eaten_at, created_at)) - ?)
				) <= 2) DESC,
				COUNT(*) DESC
		`
		err := r.db.Raw(query, userId, *hour, *hour).Scan(&rows).Error
		if err != nil {
			return nil, cerr.NewInternalError("listing meal types", err)
		}
	} else {
		query := `
			SELECT INITCAP(LOWER(MIN(type))) AS type
			FROM meals
			WHERE user_id = ?
			GROUP BY LOWER(type)
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

func applyMealFilters(q *gorm.DB, params ports.ListParams) *gorm.DB {
	if params.Date != nil {
		q = q.Where("date = ?", *params.Date)
	}
	if params.From != nil {
		q = q.Where("date >= ?", *params.From)
	}
	if params.To != nil {
		q = q.Where("date <= ?", *params.To)
	}
	if params.Type != nil {
		q = q.Where("type = ?", *params.Type)
	}
	if params.Status != nil {
		q = q.Where("status = ?", *params.Status)
	}
	if params.Tag != nil {
		q = q.Where("id IN (SELECT mtm.meal_id FROM meal_tag_map mtm JOIN meal_tags mt ON mt.id = mtm.tag_id WHERE mt.name = ?)", *params.Tag)
	}
	if params.FoodId != nil {
		q = q.Where("id IN (SELECT meal_id FROM meal_items WHERE food_id = ?)", *params.FoodId)
	}
	if params.MinCalories != nil {
		q = q.Where("calories >= ?", *params.MinCalories)
	}
	if params.MaxCalories != nil {
		q = q.Where("calories <= ?", *params.MaxCalories)
	}
	if params.MinProtein != nil {
		q = q.Where("protein_grams >= ?", *params.MinProtein)
	}
	if params.MaxProtein != nil {
		q = q.Where("protein_grams <= ?", *params.MaxProtein)
	}
	if params.MinCarbs != nil {
		q = q.Where("carbs_grams >= ?", *params.MinCarbs)
	}
	if params.MaxCarbs != nil {
		q = q.Where("carbs_grams <= ?", *params.MaxCarbs)
	}
	if params.MinFat != nil {
		q = q.Where("fat_grams >= ?", *params.MinFat)
	}
	if params.MaxFat != nil {
		q = q.Where("fat_grams <= ?", *params.MaxFat)
	}
	if params.MinFiber != nil {
		q = q.Where("fiber_grams >= ?", *params.MinFiber)
	}
	if params.MaxFiber != nil {
		q = q.Where("fiber_grams <= ?", *params.MaxFiber)
	}
	return q
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

func buildMealItem(mealId uuid.UUID, item domain.MealItem) mealItem {
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
	return mealItem{
		MealId:                 mealId,
		FoodId:                 item.FoodId,
		InputQuantity:          item.InputQuantity,
		InputUnit:              inputUnit,
		InputPortionId:         item.InputPortionId,
		InputPortionEquivalent: item.InputPortionEquivalent,
		NormalizedQuantity:     normalizedQty,
		NormalizedUnit:         normalizedUnit,
		Calories:               item.Calories,
		ProteinGrams:           item.ProteinGrams,
		CarbsGrams:             item.CarbsGrams,
		FatGrams:               item.FatGrams,
		FiberGrams:             item.FiberGrams,
		Notes:                  item.Notes,
		MeasurementMethod:      method,
	}
}

func (r *mealRepository) upsertTags(mealId, userId uuid.UUID, names []string) error {
	entries := make([]mealTag, len(names))
	for i, name := range names {
		entries[i] = mealTag{
			Id:     uuid.New(),
			UserId: userId,
			Name:   name,
		}
	}
	onConflict := clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "name"}},
		DoNothing: true,
	}
	err := r.db.Clauses(onConflict).Create(&entries).Error
	if err != nil {
		return cerr.NewInternalError("upserting meal tags", err)
	}
	var tags []mealTag
	err = r.db.Where("user_id = ? AND name IN ?", userId, names).Find(&tags).Error
	if err != nil {
		return cerr.NewInternalError("fetching meal tag ids", err)
	}
	maps := make([]mealTagMap, len(tags))
	for i, t := range tags {
		maps[i] = mealTagMap{
			MealId: mealId,
			TagId:  t.Id,
		}
	}
	mapConflict := clause.OnConflict{
		DoNothing: true,
	}
	err = r.db.Clauses(mapConflict).Create(&maps).Error
	if err != nil {
		return err
	}
	return nil
}
