package service

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"
	"github.com/ivan-ca97/life/pkg/units"
	"github.com/ivan-ca97/life/pkg/validate"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
	"github.com/ivan-ca97/life/internal/features/nutrition/ports"
)

type foodService struct {
	repository ports.FoodRepository
}

var _ ports.FoodService = (*foodService)(nil)

func NewFoodService(repository ports.FoodRepository) *foodService {
	return &foodService{
		repository: repository,
	}
}

func (s *foodService) Create(userId uuid.UUID, params ports.CreateFoodParams) (*domain.Food, error) {
	err := validate.NonEmpty(params.Name, "name")
	if err != nil {
		return nil, err
	}
	err = validateNonNegative(params.DefaultCalories, params.DefaultProteinGrams, params.DefaultCarbsGrams, params.DefaultFatGrams, params.DefaultFiberGrams)
	if err != nil {
		return nil, err
	}
	err = validateMetricUnit(params.BaseUnit)
	if err != nil {
		return nil, err
	}
	dimension, _ := units.GetDimension(params.BaseUnit)
	err = validateConversions(params.Conversions, dimension)
	if err != nil {
		return nil, err
	}

	ingredients := make([]domain.Ingredient, len(params.Ingredients))
	for i, name := range params.Ingredients {
		ingredients[i] = domain.Ingredient{
			Name: name,
		}
	}
	food := &domain.Food{
		Id:                  uuid.New(),
		UserId:              userId,
		Name:                params.Name,
		PhotoUrl:            params.PhotoUrl,
		DefaultCalories:     params.DefaultCalories,
		DefaultProteinGrams: params.DefaultProteinGrams,
		DefaultCarbsGrams:   params.DefaultCarbsGrams,
		DefaultFatGrams:     params.DefaultFatGrams,
		DefaultFiberGrams:   params.DefaultFiberGrams,
		BaseQuantity:        params.BaseQuantity,
		BaseUnit:            params.BaseUnit,
		Tags:                params.Tags,
		Ingredients:         ingredients,
		VolumeConversion:    buildVolumeConversion(params.Conversions),
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

func (s *foodService) Update(id, userId uuid.UUID, params ports.UpdateFoodParams) (*domain.Food, error) {
	if params.Name != nil {
		err := validate.NonEmpty(*params.Name, "name")
		if err != nil {
			return nil, err
		}
	}
	err := validateNonNegative(params.DefaultCalories, params.DefaultProteinGrams, params.DefaultCarbsGrams, params.DefaultFatGrams, params.DefaultFiberGrams, params.BaseQuantity)
	if err != nil {
		return nil, err
	}
	if params.BaseUnit != nil {
		err := validateMetricUnit(*params.BaseUnit)
		if err != nil {
			return nil, err
		}
	}
	if params.Conversions != nil {
		var dimension units.Dimension
		if params.BaseUnit != nil {
			dimension, _ = units.GetDimension(*params.BaseUnit)
		} else {
			current, err := s.repository.FindById(id, userId)
			if err != nil {
				return nil, err
			}
			dimension = current.Dimension()
		}
		err := validateConversions(params.Conversions, dimension)
		if err != nil {
			return nil, err
		}
	}
	food, err := s.repository.Update(id, userId, params)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *foodService) Delete(id, userId uuid.UUID) error {
	err := s.repository.Delete(id, userId)
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

func (s *foodService) Frequency(userId uuid.UUID, params ports.FrequencyParams) ([]ports.FrequencyResult, error) {
	results, err := s.repository.Frequency(userId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (s *foodService) IngredientFrequency(userId uuid.UUID, params ports.IngredientFrequencyParams) ([]ports.IngredientFrequencyResult, error) {
	results, err := s.repository.IngredientFrequency(userId, params)
	if err != nil {
		return nil, err
	}
	return results, nil
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

func (s *foodService) Impact(foodId, userId uuid.UUID) (*ports.ImpactResult, error) {
	_, err := s.repository.FindById(foodId, userId)
	if err != nil {
		return nil, err
	}
	impact, err := s.repository.Impact(foodId)
	if err != nil {
		return nil, err
	}
	return impact, nil
}

func validateMetricUnit(unit string) error {
	if !units.IsMetricUnit(unit) {
		return cerr.NewBadRequestError(fmt.Sprintf("'%s' is not a valid metric unit", unit))
	}
	return nil
}

func validateConversions(c *ports.ConversionsParam, dimension units.Dimension) error {
	if c == nil {
		return nil
	}
	if c.VolumeConversion != nil {
		if dimension == units.DimensionUnit {
			return cerr.NewBadRequestError("volume conversion not allowed for unit-type foods")
		}
		if c.VolumeConversion.GramsPerMl <= 0 {
			return cerr.NewBadRequestError("grams_per_ml must be positive")
		}
	}
	return nil
}

func validateNonNegative(values ...*float64) error {
	for _, v := range values {
		if v != nil && *v < 0 {
			return cerr.NewBadRequestError("nutritional values cannot be negative")
		}
	}
	return nil
}

func buildVolumeConversion(c *ports.ConversionsParam) *domain.VolumeConversion {
	if c == nil || c.VolumeConversion == nil {
		return nil
	}
	note := ""
	if c.VolumeConversion.Note != nil {
		note = *c.VolumeConversion.Note
	}
	return &domain.VolumeConversion{
		GramsPerMl: c.VolumeConversion.GramsPerMl,
		Note:       note,
	}
}
