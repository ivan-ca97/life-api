package handler

import (
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/nutrition/domain"
)

type volumeConversionRequest struct {
	GramsPerMl float64 `json:"grams_per_ml"`
	Note       *string `json:"note,omitempty"`
}

type conversionsRequest struct {
	VolumeConversion *volumeConversionRequest `json:"volume_conversion,omitempty"`
}

type createPortionRequest struct {
	Name           string  `json:"name"`
	BaseEquivalent float64 `json:"base_equivalent"`
}

type updatePortionRequest struct {
	Name           *string  `json:"name,omitempty"`
	BaseEquivalent *float64 `json:"base_equivalent,omitempty"`
}

type createFoodRequest struct {
	Name                string                 `json:"name"`
	PhotoUrl            string                 `json:"photo_url,omitempty"`
	DefaultCalories     *float64               `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64               `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64               `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64               `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64               `json:"default_fiber_grams,omitempty"`
	BaseQuantity        *float64               `json:"base_quantity,omitempty"`
	BaseUnit            string                 `json:"base_unit"`
	Tags                []string               `json:"tags"`
	Ingredients         []string               `json:"ingredients"`
	Conversions         *conversionsRequest    `json:"conversions,omitempty"`
	Portions            []createPortionRequest `json:"portions"`
}

type updateFoodRequest struct {
	Name                *string             `json:"name,omitempty"`
	PhotoUrl            *string             `json:"photo_url,omitempty"`
	DefaultCalories     *float64            `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64            `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64            `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64            `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64            `json:"default_fiber_grams,omitempty"`
	BaseQuantity        *float64            `json:"base_quantity,omitempty"`
	BaseUnit            *string             `json:"base_unit,omitempty"`
	Tags                *[]string           `json:"tags,omitempty"`
	Ingredients         *[]string           `json:"ingredients,omitempty"`
	Conversions         *conversionsRequest `json:"conversions,omitempty"`
}

func volumeConversionFromRequest(conversions *conversionsRequest) *domain.VolumeConversion {
	if conversions == nil || conversions.VolumeConversion == nil {
		return nil
	}
	note := ""
	if conversions.VolumeConversion.Note != nil {
		note = *conversions.VolumeConversion.Note
	}
	density := &domain.VolumeConversion{
		GramsPerMl: conversions.VolumeConversion.GramsPerMl,
		Note:       note,
	}
	return density
}

// Conversions sent without a volume conversion clear the density. Conversions
// not sent leave it unset.
func densityFromRequest(conversions *conversionsRequest) types.Optional[domain.VolumeConversion] {
	density := types.Optional[domain.VolumeConversion]{
		Set:   conversions != nil,
		Value: volumeConversionFromRequest(conversions),
	}
	return density
}
