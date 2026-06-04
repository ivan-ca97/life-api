package handler

import (
	"time"

	"github.com/google/uuid"
)

type mealItemRequest struct {
	FoodId   uuid.UUID `json:"food_id"`
	Quantity float64   `json:"quantity"`
	Unit     string    `json:"unit"`
	Notes    string    `json:"notes"`
}

type createMealRequest struct {
	Date         string            `json:"date"`
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	PhotoUrl     string            `json:"photo_url"`
	EatenAt      *time.Time        `json:"eaten_at,omitempty"`
	Calories     *float64          `json:"calories,omitempty"`
	ProteinGrams *float64          `json:"protein_grams,omitempty"`
	CarbsGrams   *float64          `json:"carbs_grams,omitempty"`
	FatGrams     *float64          `json:"fat_grams,omitempty"`
	FiberGrams   *float64          `json:"fiber_grams,omitempty"`
	Tags         []string          `json:"tags"`
	Items        []mealItemRequest `json:"items"`
	Notes        string            `json:"notes"`
}

func (r *createMealRequest) hasContent() bool {
	return r.Name != "" ||
		r.Notes != "" ||
		r.PhotoUrl != "" ||
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
	Date         *string            `json:"date,omitempty"`
	Type         *string            `json:"type,omitempty"`
	Name         *string            `json:"name,omitempty"`
	PhotoUrl     *string            `json:"photo_url,omitempty"`
	EatenAt      *time.Time         `json:"eaten_at,omitempty"`
	Calories     *float64           `json:"calories,omitempty"`
	ProteinGrams *float64           `json:"protein_grams,omitempty"`
	CarbsGrams   *float64           `json:"carbs_grams,omitempty"`
	FatGrams     *float64           `json:"fat_grams,omitempty"`
	FiberGrams   *float64           `json:"fiber_grams,omitempty"`
	Tags         *[]string          `json:"tags,omitempty"`
	Items        *[]mealItemRequest `json:"items,omitempty"`
	Notes        *string            `json:"notes,omitempty"`
}
