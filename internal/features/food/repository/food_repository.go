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
	err := r.db.Omit("Tags", "Ingredients", "Portions").Create(model).Error
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
	if len(f.Portions) > 0 {
		portions := make([]foodPortion, len(f.Portions))
		for i, p := range f.Portions {
			portions[i] = foodPortion{
				Id:             p.Id,
				FoodId:         f.Id,
				Name:           p.Name,
				BaseEquivalent: p.BaseEquivalent,
			}
		}
		err = r.db.Create(&portions).Error
		if err != nil {
			return cerr.NewInternalError("inserting food portions", err)
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
		Preload("Portions").
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

	findQuery := r.db.Preload("Tags").Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Portions").Where("user_id = ?", userId)
	if params.Query != nil {
		findQuery = findQuery.Where("name ILIKE ?", "%"+*params.Query+"%")
	}
	if params.Tag != nil {
		findQuery = findQuery.Where("id IN (SELECT food_id FROM food_tags WHERE tag = ?)", *params.Tag)
	}
	orderClause := "name ASC"
	if params.Sort != nil {
		switch *params.Sort {
		case "created_at":
			orderClause = "created_at DESC"
		case "updated_at":
			orderClause = "updated_at DESC"
		}
	}
	err = findQuery.
		Limit(params.Limit).
		Offset(params.Offset).
		Order(orderClause).
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
	if params.PhotoUrl != nil {
		updates["photo_url"] = *params.PhotoUrl
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
	if params.Public != nil {
		updates["public"] = *params.Public
	}
	if params.Conversions != nil {
		c := params.Conversions
		updates["grams_per_ml"] = nil
		updates["volume_note"] = nil
		updates["unit_base_equivalent"] = nil
		updates["unit_note"] = nil
		if c.VolumeConversion != nil {
			updates["grams_per_ml"] = c.VolumeConversion.GramsPerMl
			if c.VolumeConversion.Note != nil {
				updates["volume_note"] = *c.VolumeConversion.Note
			}
		}
		if c.UnitConversion != nil {
			updates["unit_base_equivalent"] = c.UnitConversion.BaseEquivalent
			if c.UnitConversion.Note != nil {
				updates["unit_note"] = *c.UnitConversion.Note
			}
		}
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
				tags[i] = foodTag{FoodId: id, Tag: t}
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

	if params.Portions != nil {
		if err = r.upsertPortions(id, *params.Portions); err != nil {
			return nil, err
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

func (r *foodRepository) FindByIdGlobal(id uuid.UUID) (*domain.Food, error) {
	var model food
	err := r.db.
		Preload("Tags").
		Preload("Ingredients").
		Preload("Ingredients.Ingredient").
		Preload("Portions").
		Where("id = ?", id).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrFoodNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding food by id (global)", err)
	}
	return model.toDomain(), nil
}

func (r *foodRepository) ListCommunity(params ports.CommunityListParams) (types.Page[domain.Food], error) {
	var models []food
	var total int64

	countQuery := r.db.Model(&food{}).Where("public = true")
	if params.Query != nil {
		countQuery = countQuery.Where("name ILIKE ?", "%"+*params.Query+"%")
	}
	err := countQuery.Count(&total).Error
	if err != nil {
		return types.Page[domain.Food]{}, cerr.NewInternalError("counting community foods", err)
	}

	findQuery := r.db.Preload("Tags").Preload("Ingredients").Preload("Ingredients.Ingredient").Preload("Portions").Where("public = true")
	if params.Query != nil {
		findQuery = findQuery.Where("name ILIKE ?", "%"+*params.Query+"%")
	}
	err = findQuery.
		Limit(params.Limit).
		Offset(params.Offset).
		Order("name ASC").
		Find(&models).
		Error
	if err != nil {
		return types.Page[domain.Food]{}, cerr.NewInternalError("listing community foods", err)
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

func (r *foodRepository) IsAccessibleBy(foodId, actorId uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Raw(`
		SELECT EXISTS(
			SELECT 1 FROM foods WHERE id = ? AND (public = true OR user_id = ?)
			UNION ALL
			SELECT 1 FROM foods f JOIN shares s ON s.owner_id = f.user_id AND s.resource_type = 'foods'
			WHERE f.id = ? AND s.grantee_id = ?
		)
	`, foodId, actorId, foodId, actorId).Scan(&exists).Error
	if err != nil {
		return false, cerr.NewInternalError("checking food accessibility", err)
	}
	return exists, nil
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

func (r *foodRepository) upsertPortions(foodId uuid.UUID, incoming []ports.PortionParam) error {
	var existing []foodPortion
	if err := r.db.Where("food_id = ?", foodId).Find(&existing).Error; err != nil {
		return cerr.NewInternalError("fetching food portions for upsert", err)
	}

	existingByName := make(map[string]foodPortion, len(existing))
	for _, p := range existing {
		existingByName[p.Name] = p
	}

	incomingNames := make(map[string]bool, len(incoming))
	for _, p := range incoming {
		incomingNames[p.Name] = true
	}

	for _, ex := range existing {
		if !incomingNames[ex.Name] {
			if err := r.db.Delete(&foodPortion{}, "id = ?", ex.Id).Error; err != nil {
				return cerr.NewInternalError("deleting removed food portion", err)
			}
		}
	}

	for _, p := range incoming {
		if ex, ok := existingByName[p.Name]; ok {
			if ex.BaseEquivalent != p.BaseEquivalent {
				if err := r.db.Model(&foodPortion{}).Where("id = ?", ex.Id).Update("base_equivalent", p.BaseEquivalent).Error; err != nil {
					return cerr.NewInternalError("updating food portion", err)
				}
			}
		} else {
			newPortion := foodPortion{
				Id:             uuid.New(),
				FoodId:         foodId,
				Name:           p.Name,
				BaseEquivalent: p.BaseEquivalent,
			}
			if err := r.db.Create(&newPortion).Error; err != nil {
				return cerr.NewInternalError("inserting food portion", err)
			}
		}
	}

	return nil
}

func (r *foodRepository) Impact(foodId uuid.UUID) (*ports.ImpactResult, error) {
	var totalItems int64
	if err := r.db.Raw(`SELECT COUNT(*) FROM meal_items WHERE food_id = ?`, foodId).Scan(&totalItems).Error; err != nil {
		return nil, cerr.NewInternalError("counting food impact items", err)
	}

	var totalUsers int64
	if err := r.db.Raw(`
		SELECT COUNT(DISTINCT m.user_id)
		FROM meal_items mi
		JOIN meals m ON m.id = mi.meal_id
		WHERE mi.food_id = ?
	`, foodId).Scan(&totalUsers).Error; err != nil {
		return nil, cerr.NewInternalError("counting food impact users", err)
	}

	type portionRow struct {
		PortionId   uuid.UUID `gorm:"column:portion_id"`
		PortionName string    `gorm:"column:portion_name"`
		ItemCount   int64     `gorm:"column:item_count"`
	}
	var rows []portionRow
	if err := r.db.Raw(`
		SELECT fp.id AS portion_id, fp.name AS portion_name, COUNT(mi.id) AS item_count
		FROM food_portions fp
		LEFT JOIN meal_items mi ON mi.input_portion_id = fp.id
		WHERE fp.food_id = ?
		GROUP BY fp.id, fp.name
		ORDER BY item_count DESC
	`, foodId).Scan(&rows).Error; err != nil {
		return nil, cerr.NewInternalError("querying portion impact", err)
	}

	portions := make([]ports.PortionImpact, len(rows))
	for i, row := range rows {
		portions[i] = ports.PortionImpact{
			PortionId:   row.PortionId,
			PortionName: row.PortionName,
			ItemCount:   row.ItemCount,
		}
	}

	return &ports.ImpactResult{
		TotalItems:    totalItems,
		TotalUsers:    totalUsers,
		PortionImpact: portions,
	}, nil
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
