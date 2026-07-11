package domain

// MealEstimate is the AI-produced draft the user reviews before saving. Nothing
// here is persisted by the AI: matched items reference existing foods, new-food
// suggestions are proposals the user can approve via the normal food endpoint.
type MealEstimate struct {
	MatchedItems          []MatchedItem
	NewFoodSuggestions    []NewFoodSuggestion
	Totals                Totals
	Assumptions           []string
	NeedsClarification    bool
	ClarificationQuestion string
}

// MatchedItem is a detected food the model matched to an existing catalog entry.
type MatchedItem struct {
	FoodId            string
	FoodName          string
	EstimatedQuantity float64
	Unit              string
	Confidence        string
	Assumption        string
	// SanityWarnings flags implausible stored macros for the matched food, e.g.
	// "stored calories look too low for grilled chicken".
	SanityWarnings []string
}

// NewFoodSuggestion is a food the model could not match. CreateParams is ready
// to feed the existing POST /foods endpoint if the user approves it.
type NewFoodSuggestion struct {
	Name              string
	EstimatedQuantity float64
	Unit              string
	Confidence        string
	Assumption        string
	CreateParams      SuggestedFood
}

// SuggestedFood mirrors the fields the food create endpoint needs.
type SuggestedFood struct {
	MeasurementType     string
	BaseUnit            string
	BaseQuantity        float64
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
}

type Totals struct {
	Calories     float64
	ProteinGrams float64
	CarbsGrams   float64
	FatGrams     float64
	FiberGrams   float64
}
