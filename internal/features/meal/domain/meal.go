package domain

import (
	"time"

	"github.com/google/uuid"
)

type MeasurementMethod string

const (
	MeasurementMethodPhotoEstimate MeasurementMethod = "photo_estimate"
	MeasurementMethodConfirmed     MeasurementMethod = "confirmed"
	MeasurementMethodWeighedCooked MeasurementMethod = "weighed_cooked"
	MeasurementMethodWeighedRaw    MeasurementMethod = "weighed_raw"
)

var validMeasurementMethods = map[MeasurementMethod]bool{
	MeasurementMethodPhotoEstimate: true,
	MeasurementMethodConfirmed:     true,
	MeasurementMethodWeighedCooked: true,
	MeasurementMethodWeighedRaw:    true,
}

func IsValidMeasurementMethod(m MeasurementMethod) bool {
	return m == "" || validMeasurementMethods[m]
}

type MealItem struct {
	Id                 uuid.UUID
	MealId             uuid.UUID
	FoodId             uuid.UUID
	FoodName           string
	InputQuantity      float64
	InputUnit          string
	NormalizedQuantity float64
	NormalizedUnit     string
	Calories           *float64
	ProteinGrams       *float64
	CarbsGrams         *float64
	FatGrams           *float64
	FiberGrams         *float64
	Notes              string
	MeasurementMethod  MeasurementMethod
}

type Meal struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	Date         time.Time
	Type         string
	Name         string
	PhotoUrl     string
	EatenAt      *time.Time
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
	Tags         []string
	Items        []MealItem
	Notes        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
