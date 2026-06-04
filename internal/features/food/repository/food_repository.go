package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
	"github.com/ivan-ca97/life/internal/features/food/ports"
)

type foodRepository struct {
	db *gorm.DB
}

var _ ports.FoodRepository = (*foodRepository)(nil)

func NewFoodRepository(db *gorm.DB) *foodRepository {
	return &foodRepository{
		db: db,
	}
}

func (r *foodRepository) Create(f *domain.Food) error {
	model := foodFromDomain(f)
	err := r.db.Omit("Tags", "Ingredients", "Conversions").Create(model).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return domain.ErrFoodAlreadyExists
		}
		return cerr.NewInternalError("inserting food", err)
	}
	if len(f.Tags) > 0 {
		tags := make([]foodTag, len(f.Tags))
		for i, t := range f.Tags {
			tags[i] = foodTag{
				FoodId: f.Id,
				Tag:    t,
			}
		}
		err = r.db.Create(&tags).Error
		if err != nil {
			return cerr.NewInternalError("inserting food tags", err)
		}
	}
	if len(f.Ingredients) > 0 {
		names := make([]string, len(f.Ingredients))
		for i, ing := range f.Ingredients {
			names[i] = ing.Name
		}
		ingredientIds, err := r.upsertIngredients(f.UserId, names)
		if err != nil {
			return err
		}
		for _, ingId := range ingredientIds {
			err = r.db.Exec("INSERT INTO food_ingredients (food_id, ingredient_id) VALUES (?, ?)", f.Id, ingId).Error
			if err != nil {
				return cerr.NewInternalError("inserting food ingredient link", err)
			}
		}
	}
	if len(f.Conversions) > 0 {
		conversions := make([]foodConversion, len(f.Conversions))
		for i, c := range f.Conversions {
			var note *string
			if c.Note != "" {
				n := c.Note
				note = &n
			}
			conversions[i] = foodConversion{
				Id:             c.Id,
				FoodId:         f.Id,
				Unit:           c.Unit,
				BaseEquivalent: c.BaseEquivalent,
				Inverse:        c.Inverse,
				Note:           note,
			}
		}
		err = r.db.Create(&conversions).Error
		if err != nil {
			return cerr.NewInternalError("inserting food conversions", err)
		}
	}
	return nil
}

func (r *foodRepository) FindById(id, userId uuid.UUID) (*domain.Food, error) {
	var model food
	err := r.db.
		Preload("Tags").
		Preload("Ingredients").
		Preload("Ingredients.Ingredient").
		Preload("Conversions").
		Where("id = ? AND user_id = ?", id, userId).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrFoodNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding food by id", err)
	}
	return model.toDomain(), nil
}

func (r *foodRepository) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.Food], error) {
	var models []food
	var total int64

	countQuery := r.db.Model(&food{}).Where("user_id = ?", userId)
	if params.Query != nil {
		countQuery = countQuery.Where("name ILIKE ?", "%"+*params.Query+"%")
	}
	if params.Tag != nil {
		countQuery = countQuery.Where("id IN (SELECT food_id FROM food_tags WHERE tag = ?)", *params.Tag)
	}
	err := countQuery.Count(&total).Error
	if err != nil {
		return types.Page[domain.Food]{}, cerr.NewInternalError("counting foods", err)
	}

	findQuery := r.db.Preload("Tags").Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Conversions").Where("user_id = ?", userId)
	if params.Query != nil {
		findQuery = findQuery.Where("name ILIKE ?", "%"+*params.Query+"%")
	}
	if params.Tag != nil {
		findQuery = findQuery.Where("id IN (SELECT food_id FROM food_tags WHERE tag = ?)", *params.Tag)
	}
	err = findQuery.
		Limit(params.Limit).
		Offset(params.Offset).
		Order("name ASC").
		Find(&models).
		Error
	if err != nil {
		return types.Page[domain.Food]{}, cerr.NewInternalError("listing foods", err)
	}

	foods := make([]domain.Food, len(models))
	for i, m := range models {
		foods[i] = *m.toDomain()
	}

	result := types.Page[domain.Food]{
		Items:  foods,
		Total:  total,
		Limit:  params.Limit,
		Offset: params.Offset,
	}
	return result, nil
}

func (r *foodRepository) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.Food, error) {
	var count int64
	err := r.db.Model(&food{}).Where("id = ? AND user_id = ?", id, userId).Count(&count).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking food existence", err)
	}
	if count == 0 {
		return nil, domain.ErrFoodNotFound
	}

	updates := map[string]any{}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.DefaultCalories != nil {
		updates["default_calories"] = *params.DefaultCalories
	}
	if params.DefaultProteinGrams != nil {
		updates["default_protein_grams"] = *params.DefaultProteinGrams
	}
	if params.DefaultCarbsGrams != nil {
		updates["default_carbs_grams"] = *params.DefaultCarbsGrams
	}
	if params.DefaultFatGrams != nil {
		updates["default_fat_grams"] = *params.DefaultFatGrams
	}
	if params.DefaultFiberGrams != nil {
		updates["default_fiber_grams"] = *params.DefaultFiberGrams
	}
	if params.BaseQuantity != nil {
		updates["base_quantity"] = *params.BaseQuantity
	}
	if params.BaseUnit != nil {
		updates["base_unit"] = *params.BaseUnit
	}
	if params.MeasurementType != nil {
		updates["measurement_type"] = *params.MeasurementType
	}

	if len(updates) > 0 {
		err = r.db.Model(&food{}).Where("id = ? AND user_id = ?", id, userId).Updates(updates).Error
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				return nil, domain.ErrFoodAlreadyExists
			}
			return nil, cerr.NewInternalError("updating food", err)
		}
	}

	if params.Tags != nil {
		err = r.db.Where("food_id = ?", id).Delete(&foodTag{}).Error
		if err != nil {
			return nil, cerr.NewInternalError("deleting food tags", err)
		}
		if len(*params.Tags) > 0 {
			tags := make([]foodTag, len(*params.Tags))
			for i, t := range *params.Tags {
				tags[i] = foodTag{
					FoodId: id,
					Tag:    t,
				}
			}
			err = r.db.Create(&tags).Error
			if err != nil {
				return nil, cerr.NewInternalError("inserting food tags", err)
			}
		}
	}

	if params.Ingredients != nil {
		err = r.db.Exec("DELETE FROM food_ingredients WHERE food_id = ?", id).Error
		if err != nil {
			return nil, cerr.NewInternalError("deleting food ingredients", err)
		}
		if len(*params.Ingredients) > 0 {
			ingredientIds, err := r.upsertIngredients(userId, *params.Ingredients)
			if err != nil {
				return nil, err
			}
			for _, ingId := range ingredientIds {
				err = r.db.Exec("INSERT INTO food_ingredients (food_id, ingredient_id) VALUES (?, ?)", id, ingId).Error
				if err != nil {
					return nil, cerr.NewInternalError("inserting food ingredient link", err)
				}
			}
		}
	}

	if params.Conversions != nil {
		err = r.db.Where("food_id = ?", id).Delete(&foodConversion{}).Error
		if err != nil {
			return nil, cerr.NewInternalError("deleting food conversions", err)
		}
		if len(*params.Conversions) > 0 {
			conversions := make([]foodConversion, len(*params.Conversions))
			for i, c := range *params.Conversions {
				conversions[i] = foodConversion{
					Id:             uuid.New(),
					FoodId:         id,
					Unit:           c.Unit,
					BaseEquivalent: c.BaseEquivalent,
					Inverse:        c.Inverse,
					Note:           c.Note,
				}
			}
			err = r.db.Create(&conversions).Error
			if err != nil {
				return nil, cerr.NewInternalError("inserting food conversions", err)
			}
		}
	}

	return r.FindById(id, userId)
}

func (r *foodRepository) Delete(id, userId uuid.UUID) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userId).Delete(&food{})
	if result.Error != nil {
		return cerr.NewInternalError("deleting food", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrFoodNotFound
	}
	return nil
}

func (r *foodRepository) IngredientFrequency(userId uuid.UUID, params ports.IngredientFrequencyParams) ([]ports.IngredientFrequencyResult, error) {
	query := `
		SELECT i.id AS ingredient_id, i.name AS ingredient_name, COUNT(DISTINCT m.id) AS count, MAX(m.date) AS last_date
		FROM food_ingredients fi
		JOIN ingredients i ON i.id = fi.ingredient_id
		JOIN meal_items mi ON mi.food_id = fi.food_id
		JOIN meals m ON m.id = mi.meal_id
		WHERE m.user_id = ?
	`
	args := []any{userId}

	if params.From != nil {
		query += " AND m.date >= ?"
		args = append(args, *params.From)
	}
	if params.To != nil {
		query += " AND m.date <= ?"
		args = append(args, *params.To)
	}

	query += " GROUP BY i.id, i.name ORDER BY count DESC"

	var results []ports.IngredientFrequencyResult
	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, cerr.NewInternalError("querying ingredient frequency", err)
	}
	return results, nil
}

func (r *foodRepository) ListIngredients(userId uuid.UUID, query *string) ([]domain.Ingredient, error) {
	sql := `
		SELECT i.id, i.name
		FROM ingredients i
		LEFT JOIN food_ingredients fi ON fi.ingredient_id = i.id
		WHERE i.user_id = ?
	`
	args := []any{userId}

	if query != nil && *query != "" {
		sql += " AND i.name ILIKE ?"
		args = append(args, "%"+*query+"%")
	}

	sql += " GROUP BY i.id, i.name ORDER BY COUNT(fi.food_id) DESC"

	type row struct {
		Id   uuid.UUID
		Name string
	}
	var rows []row
	err := r.db.Raw(sql, args...).Scan(&rows).Error
	if err != nil {
		return nil, cerr.NewInternalError("listing ingredients", err)
	}

	ingredients := make([]domain.Ingredient, len(rows))
	for i, r := range rows {
		ingredients[i] = domain.Ingredient{Id: r.Id, Name: r.Name}
	}
	return ingredients, nil
}

func (r *foodRepository) upsertIngredients(userId uuid.UUID, names []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(names))
	for i, name := range names {
		trimmed := strings.TrimSpace(name)
		var existing ingredient
		err := r.db.Where("user_id = ? AND lower(name) = lower(?)", userId, trimmed).First(&existing).Error
		if err == nil {
			ids[i] = existing.Id
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.NewInternalError("looking up ingredient", err)
		}
		newIng := ingredient{
			Id:     uuid.New(),
			UserId: userId,
			Name:   trimmed,
		}
		err = r.db.Create(&newIng).Error
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				err = r.db.Where("user_id = ? AND lower(name) = lower(?)", userId, trimmed).First(&existing).Error
				if err != nil {
					return nil, cerr.NewInternalError("re-fetching ingredient after conflict", err)
				}
				ids[i] = existing.Id
				continue
			}
			return nil, cerr.NewInternalError("creating ingredient", err)
		}
		ids[i] = newIng.Id
	}
	return ids, nil
}

func (r *foodRepository) Frequency(userId uuid.UUID, params ports.FrequencyParams) ([]ports.FrequencyResult, error) {
	query := `
		SELECT f.id AS food_id, f.name AS food_name, COUNT(*) AS count, MAX(m.date) AS last_date
		FROM meal_items mi
		JOIN meals m ON m.id = mi.meal_id
		JOIN foods f ON f.id = mi.food_id
		WHERE m.user_id = ?
	`
	args := []any{userId}

	if params.From != nil {
		query += " AND m.date >= ?"
		args = append(args, *params.From)
	}
	if params.To != nil {
		query += " AND m.date <= ?"
		args = append(args, *params.To)
	}
	if params.Tag != nil {
		query += " AND f.id IN (SELECT food_id FROM food_tags WHERE tag = ?)"
		args = append(args, *params.Tag)
	}

	query += " GROUP BY f.id, f.name ORDER BY count DESC"

	var results []ports.FrequencyResult
	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, cerr.NewInternalError("querying food frequency", err)
	}
	return results, nil
}
