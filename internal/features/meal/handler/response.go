package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
)

type mealPhotoResponse struct {
	Id         uuid.UUID  `json:"id"`
	Url        string     `json:"url"`
	IsPrimary  bool       `json:"is_primary"`
	MealItemId *uuid.UUID `json:"meal_item_id,omitempty"`
}

type mealItemResponse struct {
	Id                 uuid.UUID `json:"id"`
	FoodId             uuid.UUID `json:"food_id"`
	FoodName           string    `json:"food_name"`
	InputQuantity      float64   `json:"input_quantity"`
	InputUnit          string    `json:"input_unit"`
	NormalizedQuantity float64   `json:"normalized_quantity"`
	NormalizedUnit     string    `json:"normalized_unit"`
	Calories           *float64  `json:"calories,omitempty"`
	ProteinGrams       *float64  `json:"protein_grams,omitempty"`
	CarbsGrams         *float64  `json:"carbs_grams,omitempty"`
	FatGrams           *float64  `json:"fat_grams,omitempty"`
	FiberGrams         *float64  `json:"fiber_grams,omitempty"`
	Notes              string    `json:"notes"`
	MeasurementMethod  string    `json:"measurement_method,omitempty"`
}

type mealResponse struct {
	Id           uuid.UUID           `json:"id"`
	Date         string              `json:"date"`
	Type         string              `json:"type"`
	Name         string              `json:"name"`
	Status       string              `json:"status"`
	Photos       []mealPhotoResponse `json:"photos"`
	EatenAt      *time.Time          `json:"eaten_at,omitempty"`
	Calories     *float64            `json:"calories,omitempty"`
	ProteinGrams *float64            `json:"protein_grams,omitempty"`
	CarbsGrams   *float64            `json:"carbs_grams,omitempty"`
	FatGrams     *float64            `json:"fat_grams,omitempty"`
	FiberGrams   *float64            `json:"fiber_grams,omitempty"`
	Tags         []string            `json:"tags"`
	Items        []mealItemResponse  `json:"items"`
	Notes        string              `json:"notes"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

func mealFromDomain(m *domain.Meal) *mealResponse {
	tags := m.Tags
	if tags == nil {
		tags = []string{}
	}
	photos := make([]mealPhotoResponse, len(m.Photos))
	for i, p := range m.Photos {
		photos[i] = mealPhotoResponse{
			Id:         p.Id,
			Url:        p.Url,
			IsPrimary:  p.IsPrimary,
			MealItemId: p.MealItemId,
		}
	}
	items := make([]mealItemResponse, len(m.Items))
	for i, item := range m.Items {
		items[i] = mealItemResponse{
			Id:                 item.Id,
			FoodId:             item.FoodId,
			FoodName:           item.FoodName,
			InputQuantity:      item.InputQuantity,
			InputUnit:          item.InputUnit,
			NormalizedQuantity: item.NormalizedQuantity,
			NormalizedUnit:     item.NormalizedUnit,
			Calories:           item.Calories,
			ProteinGrams:       item.ProteinGrams,
			CarbsGrams:         item.CarbsGrams,
			FatGrams:           item.FatGrams,
			FiberGrams:         item.FiberGrams,
			Notes:              item.Notes,
			MeasurementMethod:  string(item.MeasurementMethod),
		}
	}
	return &mealResponse{
		Id:           m.Id,
		Date:         m.Date.Format("2006-01-02"),
		Type:         m.Type,
		Name:         m.Name,
		Status:       string(m.Status),
		Photos:       photos,
		EatenAt:      m.EatenAt,
		Calories:     m.Calories,
		ProteinGrams: m.ProteinGrams,
		CarbsGrams:   m.CarbsGrams,
		FatGrams:     m.FatGrams,
		FiberGrams:   m.FiberGrams,
		Tags:         tags,
		Items:        items,
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

type mealPage struct {
	Items  []mealResponse `json:"items"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

type nutritionPreviewItemResponse struct {
	FoodId       uuid.UUID `json:"food_id"`
	Calories     *float64  `json:"calories,omitempty"`
	ProteinGrams *float64  `json:"protein_grams,omitempty"`
	CarbsGrams   *float64  `json:"carbs_grams,omitempty"`
	FatGrams     *float64  `json:"fat_grams,omitempty"`
	FiberGrams   *float64  `json:"fiber_grams,omitempty"`
}

type nutritionPreviewResponse struct {
	Calories     *float64                       `json:"calories,omitempty"`
	ProteinGrams *float64                       `json:"protein_grams,omitempty"`
	CarbsGrams   *float64                       `json:"carbs_grams,omitempty"`
	FatGrams     *float64                       `json:"fat_grams,omitempty"`
	FiberGrams   *float64                       `json:"fiber_grams,omitempty"`
	Items        []nutritionPreviewItemResponse `json:"items"`
}

type mealTypesResponse struct {
	Types []string `json:"types"`
}

func newMealPage(page types.Page[domain.Meal]) *mealPage {
	items := make([]mealResponse, len(page.Items))
	for i, m := range page.Items {
		items[i] = *mealFromDomain(&m)
	}
	return &mealPage{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}
