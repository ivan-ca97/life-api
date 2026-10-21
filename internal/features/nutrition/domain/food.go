package domain

import (
	"time"

	"github.com/google/uuid"
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
	MeasurementType     string
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
