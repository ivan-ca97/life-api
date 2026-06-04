package domain

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	Id                  uuid.UUID
	UserId              uuid.UUID
	Name                string
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string
	BaseQuantity        float64
	BaseUnit            string
	Tags                []string
	Ingredients         []Ingredient
	Conversions         []Conversion
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Ingredient struct {
	Id   uuid.UUID
	Name string
}

type Conversion struct {
	Id             uuid.UUID
	Unit           string
	BaseEquivalent float64
	Inverse        bool
	Note           string
}
