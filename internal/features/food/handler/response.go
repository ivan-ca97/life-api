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

type volumeConversionResponse struct {
	GramsPerMl float64 `json:"grams_per_ml"`
	Note       string  `json:"note,omitempty"`
}

type unitConversionResponse struct {
	BaseEquivalent float64 `json:"base_equivalent"`
	Note           string  `json:"note,omitempty"`
}

type portionResponse struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	BaseEquivalent float64   `json:"base_equivalent"`
}

type foodResponse struct {
	Id                  uuid.UUID                 `json:"id"`
	UserId              uuid.UUID                 `json:"user_id"`
	Name                string                    `json:"name"`
	PhotoUrl            string                    `json:"photo_url,omitempty"`
	DefaultCalories     *float64                  `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64                  `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64                  `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64                  `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64                  `json:"default_fiber_grams,omitempty"`
	MeasurementType     string                    `json:"measurement_type"`
	BaseQuantity        float64                   `json:"base_quantity"`
	BaseUnit            string                    `json:"base_unit"`
	Public              bool                      `json:"public"`
	Tags                []string                  `json:"tags"`
	Ingredients         []ingredientResponse      `json:"ingredients"`
	VolumeConversion    *volumeConversionResponse `json:"volume_conversion,omitempty"`
	UnitConversion      *unitConversionResponse   `json:"unit_conversion,omitempty"`
	Portions            []portionResponse         `json:"portions"`
	CreatedAt           time.Time                 `json:"created_at"`
	UpdatedAt           time.Time                 `json:"updated_at"`
}

func foodFromDomain(f *domain.Food) *foodResponse {
	tags := f.Tags
	if tags == nil {
		tags = []string{}
	}
	ingredients := make([]ingredientResponse, len(f.Ingredients))
	for i, ing := range f.Ingredients {
		ingredients[i] = ingredientResponse{
			Id:   ing.Id,
			Name: ing.Name,
		}
	}
	portions := make([]portionResponse, len(f.Portions))
	for i, p := range f.Portions {
		portions[i] = portionResponse{
			Id:             p.Id,
			Name:           p.Name,
			BaseEquivalent: p.BaseEquivalent,
		}
	}

	var volumeConv *volumeConversionResponse
	if f.VolumeConversion != nil {
		volumeConv = &volumeConversionResponse{
			GramsPerMl: f.VolumeConversion.GramsPerMl,
			Note:       f.VolumeConversion.Note,
		}
	}
	var unitConv *unitConversionResponse
	if f.UnitConversion != nil {
		unitConv = &unitConversionResponse{
			BaseEquivalent: f.UnitConversion.BaseEquivalent,
			Note:           f.UnitConversion.Note,
		}
	}

	return &foodResponse{
		Id:                  f.Id,
		UserId:              f.UserId,
		Name:                f.Name,
		PhotoUrl:            f.PhotoUrl,
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
		VolumeConversion:    volumeConv,
		UnitConversion:      unitConv,
		Portions:            portions,
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

type portionImpactResponse struct {
	PortionId   uuid.UUID `json:"portion_id"`
	PortionName string    `json:"portion_name"`
	ItemCount   int64     `json:"item_count"`
}

type impactResponse struct {
	TotalItems    int64                   `json:"total_items"`
	TotalUsers    int64                   `json:"total_users"`
	PortionImpact []portionImpactResponse `json:"portion_impact"`
}

func newImpactResponse(r *ports.ImpactResult) *impactResponse {
	portions := make([]portionImpactResponse, len(r.PortionImpact))
	for i, p := range r.PortionImpact {
		portions[i] = portionImpactResponse{
			PortionId:   p.PortionId,
			PortionName: p.PortionName,
			ItemCount:   p.ItemCount,
		}
	}
	return &impactResponse{
		TotalItems:    r.TotalItems,
		TotalUsers:    r.TotalUsers,
		PortionImpact: portions,
	}
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
