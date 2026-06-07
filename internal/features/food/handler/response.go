package handler

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/food/domain"
	"github.com/ivan-ca97/life/internal/features/food/ports"
)

type ingredientResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ingredientsListResponse struct {
	Items []ingredientResponse `json:"items"`
}

type unitsListResponse struct {
	Mass   []string `json:"mass"`
	Volume []string `json:"volume"`
	Unit   []string `json:"unit"`
}

type foodUnitsResponse struct {
	Metric      []string `json:"metric"`
	Conversions []string `json:"conversions"`
}

type conversionResponse struct {
	Unit           string  `json:"unit"`
	BaseEquivalent float64 `json:"base_equivalent"`
	Inverse        bool    `json:"inverse"`
	Note           *string `json:"note,omitempty"`
}

type foodResponse struct {
	Id                  uuid.UUID            `json:"id"`
	UserId              uuid.UUID            `json:"user_id"`
	Name                string               `json:"name"`
	DefaultCalories     *float64             `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64             `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64             `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64             `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64             `json:"default_fiber_grams,omitempty"`
	MeasurementType     string               `json:"measurement_type"`
	BaseQuantity        float64              `json:"base_quantity"`
	BaseUnit            string               `json:"base_unit"`
	Public              bool                 `json:"public"`
	Tags                []string               `json:"tags"`
	Ingredients         []ingredientResponse   `json:"ingredients"`
	Conversions         []conversionResponse `json:"conversions"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}

func foodFromDomain(f *domain.Food) *foodResponse {
	tags := f.Tags
	if tags == nil {
		tags = []string{}
	}
	ingredients := make([]ingredientResponse, len(f.Ingredients))
	for i, ing := range f.Ingredients {
		ingredients[i] = ingredientResponse{Id: ing.Id, Name: ing.Name}
	}
	conversions := make([]conversionResponse, len(f.Conversions))
	for i, c := range f.Conversions {
		var note *string
		if c.Note != "" {
			n := c.Note
			note = &n
		}
		conversions[i] = conversionResponse{
			Unit:           c.Unit,
			BaseEquivalent: c.BaseEquivalent,
			Inverse:        c.Inverse,
			Note:           note,
		}
	}
	return &foodResponse{
		Id:                  f.Id,
		UserId:              f.UserId,
		Name:                f.Name,
		DefaultCalories:     f.DefaultCalories,
		DefaultProteinGrams: f.DefaultProteinGrams,
		DefaultCarbsGrams:   f.DefaultCarbsGrams,
		DefaultFatGrams:     f.DefaultFatGrams,
		DefaultFiberGrams:   f.DefaultFiberGrams,
		MeasurementType:     f.MeasurementType,
		BaseQuantity:        f.BaseQuantity,
		BaseUnit:            f.BaseUnit,
		Public:              f.Public,
		Tags:                tags,
		Ingredients:         ingredients,
		Conversions:         conversions,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
}

type foodPage struct {
	Items  []foodResponse `json:"items"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

func newFoodPage(page types.Page[domain.Food]) *foodPage {
	items := make([]foodResponse, len(page.Items))
	for i, f := range page.Items {
		items[i] = *foodFromDomain(&f)
	}
	return &foodPage{
		Items:  items,
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
}

type frequencyItemResponse struct {
	FoodId   uuid.UUID `json:"food_id"`
	FoodName string    `json:"food_name"`
	Count    int64     `json:"count"`
	LastDate string    `json:"last_date"`
}

type frequencyResponse struct {
	Items []frequencyItemResponse `json:"items"`
}

func newFrequencyResponse(results []ports.FrequencyResult) *frequencyResponse {
	items := make([]frequencyItemResponse, len(results))
	for i, r := range results {
		items[i] = frequencyItemResponse{
			FoodId:   r.FoodId,
			FoodName: r.FoodName,
			Count:    r.Count,
			LastDate: r.LastDate.Format("2006-01-02"),
		}
	}
	return &frequencyResponse{
		Items: items,
	}
}

type ingredientFrequencyItemResponse struct {
	IngredientId   uuid.UUID `json:"ingredient_id"`
	IngredientName string    `json:"ingredient_name"`
	Count          int64     `json:"count"`
	LastDate       string    `json:"last_date"`
}

type ingredientFrequencyResponse struct {
	Items []ingredientFrequencyItemResponse `json:"items"`
}

func newIngredientFrequencyResponse(results []ports.IngredientFrequencyResult) *ingredientFrequencyResponse {
	items := make([]ingredientFrequencyItemResponse, len(results))
	for i, r := range results {
		items[i] = ingredientFrequencyItemResponse{
			IngredientId:   r.IngredientId,
			IngredientName: r.IngredientName,
			Count:          r.Count,
			LastDate:       r.LastDate.Format("2006-01-02"),
		}
	}
	return &ingredientFrequencyResponse{
		Items: items,
	}
}
