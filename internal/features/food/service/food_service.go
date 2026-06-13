package service

import (
	"fmt"

	"github.com/google/uuid"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
	"github.com/ivan-ca97/life/internal/features/food/ports"
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

func (s *foodService) Create(userId uuid.UUID, params ports.CreateParams) (*domain.Food, error) {
	if err := validateNonNegative(params.DefaultCalories, params.DefaultProteinGrams, params.DefaultCarbsGrams, params.DefaultFatGrams, params.DefaultFiberGrams); err != nil {
		return nil, err
	}
	if err := validateMeasurementType(params.MeasurementType); err != nil {
		return nil, err
	}
	dim := domain.UnitDimension(params.MeasurementType)
	if err := validateMetricUnitForDimension(params.BaseUnit, dim); err != nil {
		return nil, err
	}
	if err := validateConversions(params.Conversions, dim); err != nil {
		return nil, err
	}
	if err := validatePortions(params.Portions); err != nil {
		return nil, err
	}

	ingredients := make([]domain.Ingredient, len(params.Ingredients))
	for i, name := range params.Ingredients {
		ingredients[i] = domain.Ingredient{Name: name}
	}
	portions := make([]domain.Portion, len(params.Portions))
	for i, p := range params.Portions {
		portions[i] = domain.Portion{
			Id:             uuid.New(),
			Name:           p.Name,
			BaseEquivalent: p.BaseEquivalent,
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
		MeasurementType:     params.MeasurementType,
		BaseQuantity:        params.BaseQuantity,
		BaseUnit:            params.BaseUnit,
		Public:              params.Public,
		Tags:                params.Tags,
		Ingredients:         ingredients,
		VolumeConversion:    buildVolumeConversion(params.Conversions),
		UnitConversion:      buildUnitConversion(params.Conversions),
		Portions:            portions,
	}
	err := s.repository.Create(food)
	if err != nil {
		return nil, err
	}
	return s.repository.FindById(food.Id, userId)
}

func (s *foodService) GetById(id, userId uuid.UUID) (*domain.Food, error) {
	food, err := s.repository.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *foodService) List(userId uuid.UUID, params ports.ListParams) (types.Page[domain.Food], error) {
	page, err := s.repository.List(userId, params)
	if err != nil {
		return types.Page[domain.Food]{}, err
	}
	return page, nil
}

func (s *foodService) Update(id, userId uuid.UUID, params ports.UpdateParams) (*domain.Food, error) {
	if err := validateNonNegative(params.DefaultCalories, params.DefaultProteinGrams, params.DefaultCarbsGrams, params.DefaultFatGrams, params.DefaultFiberGrams, params.BaseQuantity); err != nil {
		return nil, err
	}
	if params.MeasurementType != nil {
		if err := validateMeasurementType(*params.MeasurementType); err != nil {
			return nil, err
		}
	}
	if params.BaseUnit != nil {
		var dim domain.UnitDimension
		if params.MeasurementType != nil {
			dim = domain.UnitDimension(*params.MeasurementType)
		} else {
			current, err := s.repository.FindById(id, userId)
			if err != nil {
				return nil, err
			}
			dim = domain.UnitDimension(current.MeasurementType)
		}
		if err := validateMetricUnitForDimension(*params.BaseUnit, dim); err != nil {
			return nil, err
		}
	}
	if params.Conversions != nil {
		var dim domain.UnitDimension
		if params.MeasurementType != nil {
			dim = domain.UnitDimension(*params.MeasurementType)
		} else {
			current, err := s.repository.FindById(id, userId)
			if err != nil {
				return nil, err
			}
			dim = domain.UnitDimension(current.MeasurementType)
		}
		if err := validateConversions(params.Conversions, dim); err != nil {
			return nil, err
		}
	}
	if params.Portions != nil {
		if err := validatePortions(*params.Portions); err != nil {
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
	return s.repository.ListIngredients(userId, query)
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
	var unitConv *domain.UnitConversion
	if source.UnitConversion != nil {
		u := *source.UnitConversion
		unitConv = &u
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
		ingredients[i] = domain.Ingredient{Name: ing.Name}
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
		MeasurementType:     source.MeasurementType,
		BaseQuantity:        source.BaseQuantity,
		BaseUnit:            source.BaseUnit,
		Public:              false,
		Tags:                tags,
		Ingredients:         ingredients,
		VolumeConversion:    volumeConv,
		UnitConversion:      unitConv,
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

func validateMeasurementType(mt string) error {
	if !domain.IsValidMeasurementType(mt) {
		return cerr.NewBadRequestError(fmt.Sprintf("invalid measurement_type '%s', must be 'mass', 'volume' or 'unit'", mt))
	}
	return nil
}

func validateMetricUnit(unit string) error {
	if !domain.IsMetricUnit(unit) {
		return cerr.NewBadRequestError(fmt.Sprintf("'%s' is not a valid metric unit", unit))
	}
	return nil
}

func validateMetricUnitForDimension(unit string, dim domain.UnitDimension) error {
	if err := validateMetricUnit(unit); err != nil {
		return err
	}
	unitDim, _ := domain.GetUnitDimension(unit)
	if unitDim != dim {
		return cerr.NewBadRequestError(fmt.Sprintf("unit '%s' is not compatible with measurement type '%s'", unit, dim))
	}
	return nil
}

func validateConversions(c *ports.ConversionsParam, dim domain.UnitDimension) error {
	if c == nil {
		return nil
	}
	if c.VolumeConversion != nil {
		if dim == domain.DimensionUnit {
			return cerr.NewBadRequestError("volume conversion not allowed for unit-type foods")
		}
		if c.VolumeConversion.GramsPerMl <= 0 {
			return cerr.NewBadRequestError("grams_per_ml must be positive")
		}
	}
	if c.UnitConversion != nil {
		if c.UnitConversion.BaseEquivalent <= 0 {
			return cerr.NewBadRequestError("unit conversion base_equivalent must be positive")
		}
	}
	return nil
}

func validatePortions(portions []ports.PortionParam) error {
	seen := make(map[string]bool, len(portions))
	for _, p := range portions {
		if p.Name == "" {
			return cerr.NewBadRequestError("portion name is required")
		}
		if p.BaseEquivalent <= 0 {
			return cerr.NewBadRequestError(fmt.Sprintf("portion '%s' base_equivalent must be positive", p.Name))
		}
		if seen[p.Name] {
			return cerr.NewBadRequestError(fmt.Sprintf("duplicate portion name '%s'", p.Name))
		}
		seen[p.Name] = true
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

func buildUnitConversion(c *ports.ConversionsParam) *domain.UnitConversion {
	if c == nil || c.UnitConversion == nil {
		return nil
	}
	note := ""
	if c.UnitConversion.Note != nil {
		note = *c.UnitConversion.Note
	}
	return &domain.UnitConversion{
		BaseEquivalent: c.UnitConversion.BaseEquivalent,
		Note:           note,
	}
}
