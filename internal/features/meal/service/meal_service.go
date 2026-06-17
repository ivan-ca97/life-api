package service

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/dayclosure"
	"github.com/ivan-ca97/life/pkg/types"
	"github.com/ivan-ca97/life/pkg/units"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
	"github.com/ivan-ca97/life/internal/features/meal/ports"
)

const unitPorcion = "porcion"

type resolvedUnit struct {
	NormalizedQty     float64
	NormalizedUnit    string
	PortionId         *uuid.UUID
	PortionEquivalent *float64
}

type nutritionTotals struct {
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
}

type resolvedItems struct {
	Items  []domain.MealItem
	Totals nutritionTotals
}

type mealService struct {
	repository     ports.MealRepository
	foodLookup     ports.FoodLookup
	closureChecker dayclosure.DayClosureChecker
}

var _ ports.MealService = (*mealService)(nil)

func NewMealService(repository ports.MealRepository, foodLookup ports.FoodLookup, closureChecker dayclosure.DayClosureChecker) *mealService {
	return &mealService{
		repository:     repository,
		foodLookup:     foodLookup,
		closureChecker: closureChecker,
	}
}

func (s *mealService) Create(userId uuid.UUID, params ports.CreateParams) (*domain.Meal, error) {
	closed, err := s.closureChecker.IsClosed(userId, params.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	if err := validateMealParams(params.Calories, params.ProteinGrams, params.CarbsGrams, params.FatGrams, params.FiberGrams, params.Items); err != nil {
		return nil, err
	}

	resolved, err := s.resolveItems(userId, params.Items)
	if err != nil {
		return nil, err
	}

	photos, err := photosFromParams(params.Photos)
	if err != nil {
		return nil, err
	}

	meal := &domain.Meal{
		Id:           uuid.New(),
		UserId:       userId,
		Date:         params.Date,
		Type:         params.Type,
		Name:         params.Name,
		Photos:       photos,
		EatenAt:      params.EatenAt,
		Calories:     coalesce(params.Calories, resolved.Totals.Calories),
		ProteinGrams: coalesce(params.ProteinGrams, resolved.Totals.ProteinGrams),
		CarbsGrams:   coalesce(params.CarbsGrams, resolved.Totals.CarbsGrams),
		FatGrams:     coalesce(params.FatGrams, resolved.Totals.FatGrams),
		FiberGrams:   coalesce(params.FiberGrams, resolved.Totals.FiberGrams),
		Tags:         params.Tags,
		Items:        resolved.Items,
		Notes:        params.Notes,
	}
	err = s.repository.Create(meal)
	if err != nil {
		return nil, err
	}
	if len(params.Items) > 0 {
		meal, err = s.repository.FindById(meal.Id, userId)
		if err != nil {
			return nil, err
		}
	}
	return meal, nil
}

func (s *mealService) GetById(id, userId uuid.UUID) (*domain.Meal, error) {
	meal, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return meal, nil
}

func (s *mealService) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.Meal], error) {
	page, err := s.repository.List(userId, params)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	return page, nil
}

func (s *mealService) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.Meal, error) {
	meal, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	closed, err := s.closureChecker.IsClosed(userId, meal.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	var itemsForValidation []ports.ItemParam
	if params.Items != nil {
		itemsForValidation = *params.Items
	}
	if err := validateMealParams(params.Calories, params.ProteinGrams, params.CarbsGrams, params.FatGrams, params.FiberGrams, itemsForValidation); err != nil {
		return nil, err
	}
	if params.Photos != nil {
		if err := enforcePhotoParamPrimary(*params.Photos); err != nil {
			return nil, err
		}
	}
	if params.Items != nil {
		resolved, err := s.resolveItems(userId, *params.Items)
		if err != nil {
			return nil, err
		}
		params.Calories = coalesce(params.Calories, resolved.Totals.Calories)
		params.ProteinGrams = coalesce(params.ProteinGrams, resolved.Totals.ProteinGrams)
		params.CarbsGrams = coalesce(params.CarbsGrams, resolved.Totals.CarbsGrams)
		params.FatGrams = coalesce(params.FatGrams, resolved.Totals.FatGrams)
		params.FiberGrams = coalesce(params.FiberGrams, resolved.Totals.FiberGrams)
		params.ResolvedItems = &resolved.Items
	}
	updated, err := s.repository.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *mealService) PreviewNutrition(userId uuid.UUID, items []ports.ItemParam) (*ports.NutritionPreview, error) {
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, cerr.NewBadRequestError(fmt.Sprintf("item quantity must be positive for food %s", item.FoodId))
		}
	}
	resolved, err := s.resolveItems(userId, items)
	if err != nil {
		return nil, err
	}
	previewItems := make([]ports.NutritionPreviewItem, len(resolved.Items))
	for i, ri := range resolved.Items {
		previewItems[i] = ports.NutritionPreviewItem{
			FoodId:       ri.FoodId,
			Calories:     ri.Calories,
			ProteinGrams: ri.ProteinGrams,
			CarbsGrams:   ri.CarbsGrams,
			FatGrams:     ri.FatGrams,
			FiberGrams:   ri.FiberGrams,
		}
	}
	return &ports.NutritionPreview{
		Calories:     resolved.Totals.Calories,
		ProteinGrams: resolved.Totals.ProteinGrams,
		CarbsGrams:   resolved.Totals.CarbsGrams,
		FatGrams:     resolved.Totals.FatGrams,
		FiberGrams:   resolved.Totals.FiberGrams,
		Items:        previewItems,
	}, nil
}

func (s *mealService) Delete(id, userId uuid.UUID) error {
	meal, err := s.repository.FindById(id, userId)
	if err != nil {
		return err
	}
	closed, err := s.closureChecker.IsClosed(userId, meal.Date)
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

func (s *mealService) ListTypes(userId uuid.UUID, hour *int) ([]string, error) {
	types, err := s.repository.ListDistinctTypes(userId, hour)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (s *mealService) resolveItems(userId uuid.UUID, items []ports.ItemParam) (*resolvedItems, error) {
	result := &resolvedItems{}
	if len(items) == 0 {
		return result, nil
	}

	ids := make([]uuid.UUID, len(items))
	for i, item := range items {
		ids[i] = item.FoodId
	}

	foods, err := s.foodLookup.FindByIds(userId, ids)
	if err != nil {
		return nil, err
	}

	result.Items = make([]domain.MealItem, len(items))
	for i, item := range items {
		food, ok := foods[item.FoodId]
		if !ok {
			return nil, cerr.NewBadRequestError(fmt.Sprintf("food %s not found", item.FoodId))
		}

		ru, err := resolveUnit(item, food)
		if err != nil {
			return nil, err
		}

		baseInStandard, _, _ := units.ConvertToStandard(food.BaseQuantity, food.BaseUnit)
		ratio := ru.NormalizedQty / baseInStandard

		itemCalories := scale(food.DefaultCalories, ratio)
		itemProtein := scale(food.DefaultProteinGrams, ratio)
		itemCarbs := scale(food.DefaultCarbsGrams, ratio)
		itemFat := scale(food.DefaultFatGrams, ratio)
		itemFiber := scale(food.DefaultFiberGrams, ratio)

		result.Totals.Calories = addPtr(result.Totals.Calories, itemCalories)
		result.Totals.ProteinGrams = addPtr(result.Totals.ProteinGrams, itemProtein)
		result.Totals.CarbsGrams = addPtr(result.Totals.CarbsGrams, itemCarbs)
		result.Totals.FatGrams = addPtr(result.Totals.FatGrams, itemFat)
		result.Totals.FiberGrams = addPtr(result.Totals.FiberGrams, itemFiber)

		result.Items[i] = domain.MealItem{
			FoodId:                 item.FoodId,
			InputQuantity:          item.Quantity,
			InputUnit:              item.Unit,
			InputPortionId:         ru.PortionId,
			InputPortionEquivalent: ru.PortionEquivalent,
			NormalizedQuantity:     ru.NormalizedQty,
			NormalizedUnit:         ru.NormalizedUnit,
			Calories:               itemCalories,
			ProteinGrams:           itemProtein,
			CarbsGrams:             itemCarbs,
			FatGrams:               itemFat,
			FiberGrams:             itemFiber,
			Notes:                  item.Notes,
			MeasurementMethod:      item.MeasurementMethod,
		}
	}
	return result, nil
}

func resolveUnit(item ports.ItemParam, food ports.FoodNutrition) (resolvedUnit, error) {
	dim := units.Dimension(food.MeasurementType)
	standardUnit := units.StandardUnit[dim]

	// "porcion" multiplies the base serving and normalizes to the standard metric unit
	if item.Unit == unitPorcion {
		baseStd, _, err := units.ConvertToStandard(food.BaseQuantity, food.BaseUnit)
		if err != nil {
			return resolvedUnit{}, cerr.NewBadRequestError(fmt.Sprintf("food %s has invalid base unit '%s'", food.Id, food.BaseUnit))
		}
		ru := resolvedUnit{NormalizedQty: item.Quantity * baseStd, NormalizedUnit: standardUnit}
		return ru, nil
	}

	// Metric unit of the same dimension: automatic conversion
	if unitDim, ok := units.GetDimension(item.Unit); ok && unitDim == dim {
		qty, _, _ := units.ConvertToStandard(item.Quantity, item.Unit)
		ru := resolvedUnit{NormalizedQty: qty, NormalizedUnit: standardUnit}
		return ru, nil
	}

	// Cross-dimension via density (mass food + volume unit, or volume food + mass unit)
	if food.VolumeConversion != nil {
		if unitDim, ok := units.GetDimension(item.Unit); ok {
			if dim == units.DimensionMass && unitDim == units.DimensionVolume {
				mlQty, _, _ := units.ConvertToStandard(item.Quantity, item.Unit)
				ru := resolvedUnit{NormalizedQty: mlQty * food.VolumeConversion.GramsPerMl, NormalizedUnit: standardUnit}
				return ru, nil
			}
			if dim == units.DimensionVolume && unitDim == units.DimensionMass {
				gQty, _, _ := units.ConvertToStandard(item.Quantity, item.Unit)
				ru := resolvedUnit{NormalizedQty: gQty / food.VolumeConversion.GramsPerMl, NormalizedUnit: standardUnit}
				return ru, nil
			}
		}
	}

	// Named portion (captures UUID for the portion snapshot)
	for _, portion := range food.Portions {
		if item.Unit == portion.Name {
			portionId := portion.Id
			portionEq := portion.BaseEquivalent
			ru := resolvedUnit{
				NormalizedQty:     item.Quantity * portion.BaseEquivalent,
				NormalizedUnit:    standardUnit,
				PortionId:         &portionId,
				PortionEquivalent: &portionEq,
			}
			return ru, nil
		}
	}

	return resolvedUnit{}, cerr.NewBadRequestError(fmt.Sprintf("unknown unit '%s' for food %s", item.Unit, food.Id))
}

func validateMealParams(calories, protein, carbs, fat, fiber *float64, items []ports.ItemParam) error {
	for _, v := range []*float64{calories, protein, carbs, fat, fiber} {
		if v != nil && *v < 0 {
			return cerr.NewBadRequestError("nutritional values cannot be negative")
		}
	}
	seen := make(map[uuid.UUID]bool, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return cerr.NewBadRequestError(fmt.Sprintf("item quantity must be positive for food %s", item.FoodId))
		}
		if seen[item.FoodId] {
			return cerr.NewBadRequestError(fmt.Sprintf("duplicate food %s in items", item.FoodId))
		}
		seen[item.FoodId] = true
	}
	return nil
}

func scale(val *float64, ratio float64) *float64 {
	if val == nil {
		return nil
	}
	result := *val * ratio
	return &result
}

func addPtr(acc *float64, val *float64) *float64 {
	if val == nil {
		return acc
	}
	if acc == nil {
		v := *val
		return &v
	}
	sum := *acc + *val
	return &sum
}

func photosFromParams(params []ports.PhotoParam) ([]domain.MealPhoto, error) {
	photos := make([]domain.MealPhoto, len(params))
	for i, p := range params {
		photos[i] = domain.MealPhoto{
			Id:         uuid.New(),
			MealItemId: p.MealItemId,
			ItemFoodId: p.ItemFoodId,
			Url:        p.Url,
			IsPrimary:  p.IsPrimary,
		}
	}
	if err := enforcePrimary(photos); err != nil {
		return nil, err
	}
	return photos, nil
}

// enforcePrimary ensures each group (meal-level and per item) has exactly one
// primary photo. Returns an error if multiple primaries are sent for the same group.
// If a group has no primary, the first photo in that group is promoted.
// Groups are keyed by MealItemId if set, otherwise by ItemFoodId, otherwise meal-level.
func enforcePrimary(photos []domain.MealPhoto) error {
	firstPrimary := map[string]int{}
	firstInGroup := map[string]int{}

	for i, p := range photos {
		key := photoGroupKey(p)
		if _, seen := firstInGroup[key]; !seen {
			firstInGroup[key] = i
		}
		if p.IsPrimary {
			if _, has := firstPrimary[key]; !has {
				firstPrimary[key] = i
			} else {
				return cerr.NewBadRequestError("only one photo per group can be primary")
			}
		}
	}

	for key, idx := range firstInGroup {
		if _, ok := firstPrimary[key]; !ok {
			photos[idx].IsPrimary = true
		}
	}
	return nil
}

// enforcePhotoParamPrimary is the Update-path equivalent of enforcePrimary,
// operating on PhotoParam slices before they reach the repository.
func enforcePhotoParamPrimary(photos []ports.PhotoParam) error {
	firstPrimary := map[string]int{}
	firstInGroup := map[string]int{}

	for i, p := range photos {
		key := photoParamGroupKey(p)
		if _, seen := firstInGroup[key]; !seen {
			firstInGroup[key] = i
		}
		if p.IsPrimary {
			if _, has := firstPrimary[key]; !has {
				firstPrimary[key] = i
			} else {
				return cerr.NewBadRequestError("only one photo per group can be primary")
			}
		}
	}

	for key, idx := range firstInGroup {
		if _, ok := firstPrimary[key]; !ok {
			photos[idx].IsPrimary = true
		}
	}
	return nil
}

func photoParamGroupKey(p ports.PhotoParam) string {
	if p.MealItemId != nil {
		return "item:" + p.MealItemId.String()
	}
	if p.ItemFoodId != nil {
		return "food:" + p.ItemFoodId.String()
	}
	return ""
}

func photoGroupKey(p domain.MealPhoto) string {
	if p.MealItemId != nil {
		return "item:" + p.MealItemId.String()
	}
	if p.ItemFoodId != nil {
		return "food:" + p.ItemFoodId.String()
	}
	return "" // meal-level
}

func coalesce(explicit *float64, calculated *float64) *float64 {
	if explicit != nil {
		return explicit
	}
	return calculated
}
