package domain

import (
	"time"

	"github.com/google/uuid"
)

type MeasurementMethod string

const (
	MeasurementMethodWeighedRaw      MeasurementMethod = "weighed_raw"
	MeasurementMethodWeighedCooked   MeasurementMethod = "weighed_cooked"
	MeasurementMethodLabel           MeasurementMethod = "label"
	MeasurementMethodStandardPortion MeasurementMethod = "standard_portion"
	MeasurementMethodPhotoEstimate   MeasurementMethod = "photo_estimate"
	MeasurementMethodVisualEstimate  MeasurementMethod = "visual_estimate"
)

var validMeasurementMethods = map[MeasurementMethod]bool{
	MeasurementMethodWeighedRaw:      true,
	MeasurementMethodWeighedCooked:   true,
	MeasurementMethodLabel:           true,
	MeasurementMethodStandardPortion: true,
	MeasurementMethodPhotoEstimate:   true,
	MeasurementMethodVisualEstimate:  true,
}

func IsValidMeasurementMethod(m MeasurementMethod) bool {
	return m == "" || validMeasurementMethods[m]
}

type MealPhoto struct {
	Id         uuid.UUID
	MealItemId *uuid.UUID
	ItemFoodId *uuid.UUID // write-time only: resolved to MealItemId by the repository
	Url        string
	IsPrimary  bool
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
	Photos       []MealPhoto
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
