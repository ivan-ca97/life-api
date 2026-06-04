package handler

type createWeightEntryRequest struct {
	Date              string   `json:"date"`
	WeightKg          float64  `json:"weight_kg"`
	BodyFatPercentage *float64 `json:"body_fat_percentage,omitempty"`
	Notes             string   `json:"notes"`
}

type updateWeightEntryRequest struct {
	Date              *string  `json:"date,omitempty"`
	WeightKg          *float64 `json:"weight_kg,omitempty"`
	BodyFatPercentage *float64 `json:"body_fat_percentage,omitempty"`
	Notes             *string  `json:"notes,omitempty"`
}
