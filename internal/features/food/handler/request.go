package handler

type volumeConversionRequest struct {
	GramsPerMl float64 `json:"grams_per_ml"`
	Note       *string `json:"note,omitempty"`
}

type unitConversionRequest struct {
	BaseEquivalent float64 `json:"base_equivalent"`
	Note           *string `json:"note,omitempty"`
}

type conversionsRequest struct {
	VolumeConversion *volumeConversionRequest `json:"volume_conversion,omitempty"`
	UnitConversion   *unitConversionRequest   `json:"unit_conversion,omitempty"`
}

type portionRequest struct {
	Name           string  `json:"name"`
	BaseEquivalent float64 `json:"base_equivalent"`
}

type createFoodRequest struct {
	Name                string              `json:"name"`
	DefaultCalories     *float64            `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64            `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64            `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64            `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64            `json:"default_fiber_grams,omitempty"`
	MeasurementType     string              `json:"measurement_type"`
	BaseQuantity        *float64            `json:"base_quantity,omitempty"`
	BaseUnit            string              `json:"base_unit"`
	Public              bool                `json:"public"`
	Tags                []string            `json:"tags"`
	Ingredients         []string            `json:"ingredients"`
	Conversions         *conversionsRequest `json:"conversions,omitempty"`
	Portions            []portionRequest    `json:"portions"`
}

type updateFoodRequest struct {
	Name                *string              `json:"name,omitempty"`
	DefaultCalories     *float64             `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64             `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64             `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64             `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64             `json:"default_fiber_grams,omitempty"`
	MeasurementType     *string              `json:"measurement_type,omitempty"`
	BaseQuantity        *float64             `json:"base_quantity,omitempty"`
	BaseUnit            *string              `json:"base_unit,omitempty"`
	Public              *bool                `json:"public,omitempty"`
	Tags                *[]string            `json:"tags,omitempty"`
	Ingredients         *[]string            `json:"ingredients,omitempty"`
	Conversions         *conversionsRequest  `json:"conversions,omitempty"`
	Portions            *[]portionRequest    `json:"portions,omitempty"`
}
