package handler

type conversionRequest struct {
	Unit           string  `json:"unit"`
	BaseEquivalent float64 `json:"base_equivalent"`
	Inverse        bool    `json:"inverse"`
	Note           *string `json:"note,omitempty"`
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
	Tags                []string            `json:"tags"`
	Ingredients         []string            `json:"ingredients"`
	Conversions         []conversionRequest `json:"conversions"`
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
	Tags                *[]string            `json:"tags,omitempty"`
	Ingredients         *[]string            `json:"ingredients,omitempty"`
	Conversions         *[]conversionRequest `json:"conversions,omitempty"`
}
