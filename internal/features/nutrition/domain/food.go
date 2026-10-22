package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/units"
)

type VolumeConversion struct {
	GramsPerMl float64
	Note       string
}

type Portion struct {
	Id             uuid.UUID
	Name           string
	BaseEquivalent float64
}

type Food struct {
	Id                  uuid.UUID
	UserId              uuid.UUID
	Name                string
	PhotoUrl            string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	BaseQuantity        float64
	BaseUnit            units.Unit
	Tags                []string
	Ingredients         []Ingredient
	VolumeConversion    *VolumeConversion
	Portions            []Portion
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Ingredient struct {
	Id   uuid.UUID
	Name string
}

// Comes from the base unit, which is always a metric unit.
func (f Food) Dimension() units.Dimension {
	return f.BaseUnit.Dimension()
}

func (f Food) StandardUnit() string {
	return units.StandardUnit[f.Dimension()]
}

// Converts a quantity in a metric unit to the food's standard unit: directly for
// a unit of the food's own dimension, through its density for one of the other.
func (f Food) InStandardUnit(quantity float64, unit string) (float64, error) {
	metricUnit, err := units.ParseUnit(unit)
	if err != nil {
		return 0, NewUnknownUnitError(f.Id, unit)
	}
	inStandardUnit := metricUnit.ToStandard(quantity)

	switch {
	case metricUnit.Dimension() == f.Dimension():
		return inStandardUnit, nil
	case f.canConvertByDensity(metricUnit.Dimension()):
		return f.convertByDensity(inStandardUnit), nil
	default:
		return 0, NewUnknownUnitError(f.Id, unit)
	}
}

// Mass and volume convert into each other when the food has a density.
func (f Food) canConvertByDensity(dimension units.Dimension) bool {
	if f.VolumeConversion == nil {
		return false
	}
	massToVolume := f.Dimension() == units.DimensionVolume && dimension == units.DimensionMass
	volumeToMass := f.Dimension() == units.DimensionMass && dimension == units.DimensionVolume
	return massToVolume || volumeToMass
}

// Takes a quantity of the other dimension, in its standard unit, to the food's
// standard unit.
func (f Food) convertByDensity(quantity float64) float64 {
	gramsPerMl := f.VolumeConversion.GramsPerMl
	if f.Dimension() == units.DimensionMass {
		return quantity * gramsPerMl
	}
	return quantity / gramsPerMl
}

func (f Food) Portion(portionId uuid.UUID) (*Portion, error) {
	for _, portion := range f.Portions {
		if portion.Id == portionId {
			found := portion
			return &found, nil
		}
	}
	return nil, NewUnknownPortionError(f.Id, portionId)
}

// Scales the nutrients, given for the base quantity, to a quantity in the
// standard unit. An undefined nutrient stays undefined.
func (f Food) NutritionFor(standardQuantity float64) NutritionTotals {
	ratio := standardQuantity / f.BaseUnit.ToStandard(f.BaseQuantity)
	nutrition := NutritionTotals{
		Calories:     scaleNutrient(f.DefaultCalories, ratio),
		ProteinGrams: scaleNutrient(f.DefaultProteinGrams, ratio),
		CarbsGrams:   scaleNutrient(f.DefaultCarbsGrams, ratio),
		FatGrams:     scaleNutrient(f.DefaultFatGrams, ratio),
		FiberGrams:   scaleNutrient(f.DefaultFiberGrams, ratio),
	}
	return nutrition
}

func scaleNutrient(amount *float64, ratio float64) *float64 {
	if amount == nil {
		return nil
	}
	scaled := *amount * ratio
	return &scaled
}

// Builds a food whose nutrients refer to baseQuantity of a metric base unit.
func NewFood(userId uuid.UUID, name string, baseQuantity float64, baseUnit string) (*Food, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrFoodNameRequired
	}
	unit, err := ParseBaseUnit(baseUnit)
	if err != nil {
		return nil, err
	}
	food := &Food{
		Id:           uuid.New(),
		UserId:       userId,
		Name:         name,
		BaseQuantity: baseQuantity,
		BaseUnit:     *unit,
	}
	return food, nil
}

// A food's base unit is always a metric unit.
func ParseBaseUnit(name string) (*units.Unit, error) {
	unit, err := units.ParseUnit(name)
	if err != nil {
		return nil, NewInvalidBaseUnitError(name)
	}
	return unit, nil
}

// Builds a valid portion: the name is required and the equivalent, in the food's
// base unit, is positive.
func NewPortion(name string, baseEquivalent float64) (*Portion, error) {
	err := validatePortion(name, baseEquivalent)
	if err != nil {
		return nil, err
	}
	portion := &Portion{
		Id:             uuid.New(),
		Name:           name,
		BaseEquivalent: baseEquivalent,
	}
	return portion, nil
}

func validatePortion(name string, baseEquivalent float64) error {
	if name == "" {
		return ErrPortionNameRequired
	}
	if baseEquivalent <= 0 {
		return NewInvalidPortionEquivalentError(name)
	}
	return nil
}

func (f *Food) AddPortion(name string, baseEquivalent float64) (*Portion, error) {
	portion, err := NewPortion(name, baseEquivalent)
	if err != nil {
		return nil, err
	}
	if f.hasPortionNamed(name, portion.Id) {
		return nil, ErrPortionAlreadyExists
	}
	f.Portions = append(f.Portions, *portion)
	return portion, nil
}

func (f *Food) UpdatePortion(portionId uuid.UUID, name *string, baseEquivalent *float64) (*Portion, error) {
	index := -1
	for i, portion := range f.Portions {
		if portion.Id == portionId {
			index = i
		}
	}
	if index < 0 {
		return nil, ErrPortionNotFound
	}

	updated := f.Portions[index]
	if name != nil {
		updated.Name = *name
	}
	if baseEquivalent != nil {
		updated.BaseEquivalent = *baseEquivalent
	}
	err := validatePortion(updated.Name, updated.BaseEquivalent)
	if err != nil {
		return nil, err
	}
	if f.hasPortionNamed(updated.Name, portionId) {
		return nil, ErrPortionAlreadyExists
	}
	f.Portions[index] = updated
	return &updated, nil
}

func (f *Food) hasPortionNamed(name string, exceptId uuid.UUID) bool {
	for _, portion := range f.Portions {
		if portion.Name == name && portion.Id != exceptId {
			return true
		}
	}
	return false
}
