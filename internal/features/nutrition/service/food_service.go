package service

import (
	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
	"github.com/ivan-ca97/life/internal/features/nutrition/domain/food_change"
	"github.com/ivan-ca97/life/internal/features/nutrition/domain/food_option"
	"github.com/ivan-ca97/life/internal/features/nutrition/ports"
)

type foodService struct {
	repository ports.FoodRepository
	usage      ports.FoodUsage
}

var _ ports.FoodService = (*foodService)(nil)

func NewFoodService(repository ports.FoodRepository, usage ports.FoodUsage) *foodService {
	return &foodService{
		repository: repository,
		usage:      usage,
	}
}

func (s *foodService) Create(userId uuid.UUID, params ports.CreateFoodParams) (*domain.Food, error) {
	food, err := domain.NewFood(
		userId,
		params.Name,
		params.BaseQuantity,
		params.BaseUnit,
		food_option.Photo(&params.PhotoUrl),
		food_option.Calories(params.DefaultCalories),
		food_option.ProteinGrams(params.DefaultProteinGrams),
		food_option.CarbsGrams(params.DefaultCarbsGrams),
		food_option.FatGrams(params.DefaultFatGrams),
		food_option.FiberGrams(params.DefaultFiberGrams),
		food_option.Density(params.Density),
		food_option.Tags(&params.Tags),
		food_option.Ingredients(&params.Ingredients),
	)
	if err != nil {
		return nil, err
	}
	for _, portion := range params.Portions {
		_, err = food.AddPortion(portion.Name, portion.BaseEquivalent)
		if err != nil {
			return nil, err
		}
	}
	err = s.repository.Create(food)
	if err != nil {
		return nil, err
	}
	created, err := s.repository.FindById(food.Id, userId)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *foodService) GetById(id, userId uuid.UUID) (*domain.Food, error) {
	food, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return food, nil
}

// FindUsableByIds returns the foods among ids that userId can use.
// IDs the user cannot use, or that do not exist, are left out.
func (s *foodService) FindUsableByIds(userId uuid.UUID, ids []uuid.UUID) ([]domain.Food, error) {
	foods, err := s.repository.FindUsableByIds(userId, ids)
	if err != nil {
		return nil, err
	}
	return foods, nil
}

func (s *foodService) List(userId uuid.UUID, params ports.ListFoodParams) (types.Page[domain.Food], error) {
	page, err := s.repository.List(userId, params)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	return page, nil
}

// A field missing from the request keeps its current value.
func (s *foodService) Update(id, userId uuid.UUID, params ports.UpdateFoodParams) (*domain.Food, error) {
	food, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}

	err = food.Update(
		food_change.Name(params.Name),
		food_change.Photo(params.PhotoUrl),
		food_change.Calories(params.DefaultCalories),
		food_change.ProteinGrams(params.DefaultProteinGrams),
		food_change.CarbsGrams(params.DefaultCarbsGrams),
		food_change.FatGrams(params.DefaultFatGrams),
		food_change.FiberGrams(params.DefaultFiberGrams),
		food_change.BaseQuantity(params.BaseQuantity),
		food_change.BaseUnit(params.BaseUnit),
		food_change.Density(params.Density),
		food_change.Tags(params.Tags),
		food_change.Ingredients(params.Ingredients),
	)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(food)
	if err != nil {
		return nil, err
	}
	updated, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// A food that meals were logged with stays, so those meals keep their food.
func (s *foodService) Delete(id, userId uuid.UUID) error {
	_, err := s.repository.FindById(id, userId)
	if err != nil {
		return err
	}
	count, err := s.usage.CountMealItems(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return domain.ErrFoodInUse
	}
	err = s.repository.Delete(id, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *foodService) AddPortion(foodId, userId uuid.UUID, params ports.PortionParam) (*domain.Portion, error) {
	food, err := s.repository.FindById(foodId, userId)
	if err != nil {
		return nil, err
	}
	portion, err := food.AddPortion(params.Name, params.BaseEquivalent)
	if err != nil {
		return nil, err
	}
	err = s.repository.AddPortion(foodId, userId, portion)
	if err != nil {
		return nil, err
	}
	return portion, nil
}

func (s *foodService) UpdatePortion(foodId, portionId, userId uuid.UUID, params ports.UpdatePortionParams) (*domain.Portion, error) {
	food, err := s.repository.FindById(foodId, userId)
	if err != nil {
		return nil, err
	}
	portion, err := food.UpdatePortion(portionId, params.Name, params.BaseEquivalent)
	if err != nil {
		return nil, err
	}
	err = s.repository.UpdatePortion(foodId, userId, portion)
	if err != nil {
		return nil, err
	}
	return portion, nil
}

func (s *foodService) DeletePortion(foodId, portionId, userId uuid.UUID) error {
	err := s.repository.DeletePortion(foodId, portionId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *foodService) ListIngredients(userId uuid.UUID, query *string) ([]domain.Ingredient, error) {
	ingredients, err := s.repository.ListIngredients(userId, query)
	if err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (s *foodService) ListCommunity(params ports.CommunityListParams) (types.Page[domain.Food], error) {
	page, err := s.repository.ListCommunity(params)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	return page, nil
}

func (s *foodService) Copy(actorId, foodId uuid.UUID) (*domain.Food, error) {
	accessible, err := s.repository.IsAccessibleBy(foodId, actorId)
	if err != nil {
		return nil, err
	}
	if !accessible {
		return nil, cerr.NewForbiddenError("you do not have access to this food")
	}

	source, err := s.repository.FindByIdGlobal(foodId)
	if err != nil {
		return nil, err
	}

	var volumeConv *domain.VolumeConversion
	if source.VolumeConversion != nil {
		v := *source.VolumeConversion
		volumeConv = &v
	}
	portions := make([]domain.Portion, len(source.Portions))
	for i, p := range source.Portions {
		portions[i] = domain.Portion{
			Id:             uuid.New(),
			Name:           p.Name,
			BaseEquivalent: p.BaseEquivalent,
		}
	}
	ingredients := make([]domain.Ingredient, len(source.Ingredients))
	for i, ing := range source.Ingredients {
		ingredients[i] = domain.Ingredient{
			Name: ing.Name,
		}
	}
	tags := make([]string, len(source.Tags))
	copy(tags, source.Tags)

	copied := &domain.Food{
		Id:                  uuid.New(),
		UserId:              actorId,
		Name:                source.Name,
		DefaultCalories:     source.DefaultCalories,
		DefaultProteinGrams: source.DefaultProteinGrams,
		DefaultCarbsGrams:   source.DefaultCarbsGrams,
		DefaultFatGrams:     source.DefaultFatGrams,
		DefaultFiberGrams:   source.DefaultFiberGrams,
		BaseQuantity:        source.BaseQuantity,
		BaseUnit:            source.BaseUnit,
		Tags:                tags,
		Ingredients:         ingredients,
		VolumeConversion:    volumeConv,
		Portions:            portions,
	}

	err = s.repository.Create(copied)
	if err != nil {
		return nil, err
	}

	result, err := s.repository.FindById(copied.Id, actorId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
