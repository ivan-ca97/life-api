package service

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/dayclosure"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
	"github.com/ivan-ca97/life/internal/features/nutrition/ports"
)

type mealService struct {
	repository     ports.MealRepository
	foodService    ports.FoodService
	closureChecker dayclosure.DayClosureChecker
}

var _ ports.MealService = (*mealService)(nil)

func NewMealService(repository ports.MealRepository, foodService ports.FoodService, closureChecker dayclosure.DayClosureChecker) *mealService {
	return &mealService{
		repository:     repository,
		foodService:    foodService,
		closureChecker: closureChecker,
	}
}

func (s *mealService) Create(userId uuid.UUID, params ports.CreateMealParams) (*domain.Meal, error) {
	closed, err := s.closureChecker.IsClosed(userId, params.Date)
	if err != nil {
		return nil, err
	}
	if closed {
		return nil, dayclosure.ErrDayClosed
	}

	err = validateMealParams(params.Calories, params.ProteinGrams, params.CarbsGrams, params.FatGrams, params.FiberGrams, params.Items)
	if err != nil {
		return nil, err
	}

	status := domain.MealStatus(params.Status)
	if status == "" {
		status = domain.MealStatusComplete
	} else if !domain.IsValidMealStatus(status) {
		return nil, cerr.NewBadRequestError("status must be 'complete' or 'pending'")
	}

	items, err := s.resolveItems(userId, params.Items)
	if err != nil {
		return nil, err
	}
	totals := domain.TotalNutrition(items)

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
		Status:       status,
		Photos:       photos,
		EatenAt:      params.EatenAt,
		Calories:     coalesce(params.Calories, totals.Calories),
		ProteinGrams: coalesce(params.ProteinGrams, totals.ProteinGrams),
		CarbsGrams:   coalesce(params.CarbsGrams, totals.CarbsGrams),
		FatGrams:     coalesce(params.FatGrams, totals.FatGrams),
		FiberGrams:   coalesce(params.FiberGrams, totals.FiberGrams),
		Tags:         params.Tags,
		Items:        items,
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

func (s *mealService) List(userId uuid.UUID, params ports.ListMealParams) (types.Page[domain.Meal], error) {
	page, err := s.repository.List(userId, params)
	if err != nil {
		return types.Page[domain.Meal]{}, err
	}
	return page, nil
}

func (s *mealService) Update(id, userId uuid.UUID, params ports.UpdateMealParams) (*domain.Meal, error) {
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
	err = validateMealParams(params.Calories, params.ProteinGrams, params.CarbsGrams, params.FatGrams, params.FiberGrams, itemsForValidation)
	if err != nil {
		return nil, err
	}
	if params.Status != nil {
		if !domain.IsValidMealStatus(domain.MealStatus(*params.Status)) {
			return nil, cerr.NewBadRequestError("status must be 'complete' or 'pending'")
		}
	}
	if params.Photos != nil {
		err := enforcePhotoParamPrimary(*params.Photos)
		if err != nil {
			return nil, err
		}
	}
	if params.Items != nil {
		items, err := s.resolveItems(userId, *params.Items)
		if err != nil {
			return nil, err
		}
		totals := domain.TotalNutrition(items)
		params.Calories = coalesce(params.Calories, totals.Calories)
		params.ProteinGrams = coalesce(params.ProteinGrams, totals.ProteinGrams)
		params.CarbsGrams = coalesce(params.CarbsGrams, totals.CarbsGrams)
		params.FatGrams = coalesce(params.FatGrams, totals.FatGrams)
		params.FiberGrams = coalesce(params.FiberGrams, totals.FiberGrams)
		params.ResolvedItems = &items
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
	totals := domain.TotalNutrition(resolved)
	previewItems := make([]ports.NutritionPreviewItem, len(resolved))
	for i, item := range resolved {
		previewItems[i] = ports.NutritionPreviewItem{
			FoodId:       item.FoodId,
			Calories:     item.Calories,
			ProteinGrams: item.ProteinGrams,
			CarbsGrams:   item.CarbsGrams,
			FatGrams:     item.FatGrams,
			FiberGrams:   item.FiberGrams,
		}
	}
	result := &ports.NutritionPreview{
		Calories:     totals.Calories,
		ProteinGrams: totals.ProteinGrams,
		CarbsGrams:   totals.CarbsGrams,
		FatGrams:     totals.FatGrams,
		FiberGrams:   totals.FiberGrams,
		Items:        previewItems,
	}
	return result, nil
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

func (s *mealService) resolveItems(userId uuid.UUID, items []ports.ItemParam) ([]domain.MealItem, error) {
	if len(items) == 0 {
		return nil, nil
	}

	ids := make([]uuid.UUID, len(items))
	for i, item := range items {
		ids[i] = item.FoodId
	}

	usableFoods, err := s.foodService.FindUsableByIds(userId, ids)
	if err != nil {
		return nil, err
	}
	foods := make(map[uuid.UUID]domain.Food, len(usableFoods))
	for _, food := range usableFoods {
		foods[food.Id] = food
	}

	resolved := make([]domain.MealItem, len(items))
	for i, item := range items {
		food, ok := foods[item.FoodId]
		if !ok {
			return nil, cerr.NewBadRequestError(fmt.Sprintf("food %s not found", item.FoodId))
		}

		mealItem, err := newMealItem(food, item)
		if err != nil {
			return nil, err
		}
		mealItem.Notes = item.Notes
		mealItem.MeasurementMethod = item.MeasurementMethod
		resolved[i] = *mealItem
	}
	return resolved, nil
}

func newMealItem(food domain.Food, item ports.ItemParam) (*domain.MealItem, error) {
	if item.PortionId == nil {
		mealItem, err := domain.NewMealItem(food, item.Quantity, item.Unit)
		if err != nil {
			return nil, err
		}
		return mealItem, nil
	}
	if item.Unit != "" {
		return nil, domain.ErrItemUnitAndPortion
	}
	mealItem, err := domain.NewMealItemOfPortion(food, item.Quantity, *item.PortionId)
	if err != nil {
		return nil, err
	}
	return mealItem, nil
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
	err := enforcePrimary(photos)
	if err != nil {
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
		_, seen := firstInGroup[key]
		if !seen {
			firstInGroup[key] = i
		}
		if p.IsPrimary {
			_, has := firstPrimary[key]
			if !has {
				firstPrimary[key] = i
			} else {
				return cerr.NewBadRequestError("only one photo per group can be primary")
			}
		}
	}

	for key, idx := range firstInGroup {
		_, ok := firstPrimary[key]
		if !ok {
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
		_, seen := firstInGroup[key]
		if !seen {
			firstInGroup[key] = i
		}
		if p.IsPrimary {
			_, has := firstPrimary[key]
			if !has {
				firstPrimary[key] = i
			} else {
				return cerr.NewBadRequestError("only one photo per group can be primary")
			}
		}
	}

	for key, idx := range firstInGroup {
		_, ok := firstPrimary[key]
		if !ok {
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
