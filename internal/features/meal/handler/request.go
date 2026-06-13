package handler

import (
	"time"

	"github.com/google/uuid"
)

type mealPhotoRequest struct {
	Url        string     `json:"url"`
	IsPrimary  bool       `json:"is_primary"`
	MealItemId *uuid.UUID `json:"meal_item_id,omitempty"`
}

type mealItemRequest struct {
	FoodId            uuid.UUID `json:"food_id"`
	Quantity          float64   `json:"quantity"`
	Unit              string    `json:"unit"`
	Notes             string    `json:"notes"`
	MeasurementMethod string    `json:"measurement_method,omitempty"`
}

type createMealRequest struct {
	Date         string             `json:"date"`
	Type         string             `json:"type"`
	Name         string             `json:"name"`
	Photos       []mealPhotoRequest `json:"photos,omitempty"`
	EatenAt      *time.Time         `json:"eaten_at,omitempty"`
	Calories     *float64           `json:"calories,omitempty"`
	ProteinGrams *float64           `json:"protein_grams,omitempty"`
	CarbsGrams   *float64           `json:"carbs_grams,omitempty"`
	FatGrams     *float64           `json:"fat_grams,omitempty"`
	FiberGrams   *float64           `json:"fiber_grams,omitempty"`
	Tags         []string           `json:"tags"`
	Items        []mealItemRequest  `json:"items"`
	Notes        string             `json:"notes"`
}

func (r *createMealRequest) hasContent() bool {
	return r.Name != "" ||
		r.Notes != "" ||
		len(r.Photos) > 0 ||
		r.Calories != nil ||
		r.ProteinGrams != nil ||
		r.CarbsGrams != nil ||
		r.FatGrams != nil ||
		r.FiberGrams != nil ||
		len(r.Items) > 0
}

type previewNutritionRequest struct {
	Items []mealItemRequest `json:"items"`
}

type updateMealRequest struct {
	Date         *string             `json:"date,omitempty"`
	Type         *string             `json:"type,omitempty"`
	Name         *string             `json:"name,omitempty"`
	Photos       *[]mealPhotoRequest `json:"photos,omitempty"`
	EatenAt      *time.Time          `json:"eaten_at,omitempty"`
	Calories     *float64            `json:"calories,omitempty"`
	ProteinGrams *float64            `json:"protein_grams,omitempty"`
	CarbsGrams   *float64            `json:"carbs_grams,omitempty"`
	FatGrams     *float64            `json:"fat_grams,omitempty"`
	FiberGrams   *float64            `json:"fiber_grams,omitempty"`
	Tags         *[]string           `json:"tags,omitempty"`
	Items        *[]mealItemRequest  `json:"items,omitempty"`
	Notes        *string             `json:"notes,omitempty"`
}
