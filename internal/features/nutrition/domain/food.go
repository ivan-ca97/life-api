package domain

import (
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
	BaseUnit            string
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
	dimension, _ := units.GetDimension(f.BaseUnit)
	return dimension
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
