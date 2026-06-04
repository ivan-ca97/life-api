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
	if err := validateConversions(params.Conversions); err != nil {
		return nil, err
	}
	conversions := make([]domain.Conversion, len(params.Conversions))
	for i, c := range params.Conversions {
		note := ""
		if c.Note != nil {
			note = *c.Note
		}
		conversions[i] = domain.Conversion{
			Id:             uuid.New(),
			Unit:           c.Unit,
			BaseEquivalent: c.BaseEquivalent,
			Inverse:        c.Inverse,
			Note:           note,
		}
	}
	ingredients := make([]domain.Ingredient, len(params.Ingredients))
	for i, name := range params.Ingredients {
		ingredients[i] = domain.Ingredient{Name: name}
	}
	food := &domain.Food{
		Id:                  uuid.New(),
		UserId:              userId,
		Name:                params.Name,
		DefaultCalories:     params.DefaultCalories,
		DefaultProteinGrams: params.DefaultProteinGrams,
		DefaultCarbsGrams:   params.DefaultCarbsGrams,
		DefaultFatGrams:     params.DefaultFatGrams,
		DefaultFiberGrams:   params.DefaultFiberGrams,
		MeasurementType:     params.MeasurementType,
		BaseQuantity:        params.BaseQuantity,
		BaseUnit:            params.BaseUnit,
		Tags:                params.Tags,
		Ingredients:         ingredients,
		Conversions:         conversions,
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
		if err := validateConversions(*params.Conversions); err != nil {
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

func validateConversions(conversions []ports.ConversionParam) error {
	seen := make(map[string]bool, len(conversions))
	for _, c := range conversions {
		if c.Unit == "" {
			return cerr.NewBadRequestError("conversion unit is required")
		}
		if c.BaseEquivalent <= 0 {
			return cerr.NewBadRequestError(fmt.Sprintf("base_equivalent must be positive for conversion '%s'", c.Unit))
		}
		if seen[c.Unit] {
			return cerr.NewBadRequestError(fmt.Sprintf("duplicate conversion unit '%s'", c.Unit))
		}
		seen[c.Unit] = true
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
