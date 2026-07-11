package handler

import "github.com/ivan-ca97/life/internal/applications/meal_ai/domain"

type matchedItemResponse struct {
	FoodId            string   `json:"food_id"`
	FoodName          string   `json:"food_name"`
	EstimatedQuantity float64  `json:"estimated_quantity"`
	Unit              string   `json:"unit"`
	Confidence        string   `json:"confidence"`
	Assumption        string   `json:"assumption"`
	SanityWarnings    []string `json:"sanity_warnings"`
}

type suggestedFoodResponse struct {
	MeasurementType     string   `json:"measurement_type"`
	BaseUnit            string   `json:"base_unit"`
	BaseQuantity        float64  `json:"base_quantity"`
	DefaultCalories     *float64 `json:"default_calories"`
	DefaultProteinGrams *float64 `json:"default_protein_grams"`
	DefaultCarbsGrams   *float64 `json:"default_carbs_grams"`
	DefaultFatGrams     *float64 `json:"default_fat_grams"`
	DefaultFiberGrams   *float64 `json:"default_fiber_grams"`
}

type newFoodSuggestionResponse struct {
	Name              string                `json:"name"`
	EstimatedQuantity float64               `json:"estimated_quantity"`
	Unit              string                `json:"unit"`
	Confidence        string                `json:"confidence"`
	Assumption        string                `json:"assumption"`
	CreateParams      suggestedFoodResponse `json:"create_params"`
}

type totalsResponse struct {
	Calories     float64 `json:"calories"`
	ProteinGrams float64 `json:"protein_grams"`
	CarbsGrams   float64 `json:"carbs_grams"`
	FatGrams     float64 `json:"fat_grams"`
	FiberGrams   float64 `json:"fiber_grams"`
}

type usageResponse struct {
	Model        string  `json:"model"`
	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CostUsd      float64 `json:"cost_usd"`
}

type estimateResponse struct {
	MatchedItems          []matchedItemResponse       `json:"matched_items"`
	NewFoodSuggestions    []newFoodSuggestionResponse `json:"new_food_suggestions"`
	Totals                totalsResponse              `json:"totals"`
	Assumptions           []string                    `json:"assumptions"`
	NeedsClarification    bool                        `json:"needs_clarification"`
	ClarificationQuestion string                      `json:"clarification_question,omitempty"`
	Usage                 usageResponse               `json:"usage"`
}

func estimateFromDomain(e *domain.MealEstimate) *estimateResponse {
	resp := &estimateResponse{
		Totals: totalsResponse{
			Calories:     e.Totals.Calories,
			ProteinGrams: e.Totals.ProteinGrams,
			CarbsGrams:   e.Totals.CarbsGrams,
			FatGrams:     e.Totals.FatGrams,
			FiberGrams:   e.Totals.FiberGrams,
		},
		Assumptions:           e.Assumptions,
		NeedsClarification:    e.NeedsClarification,
		ClarificationQuestion: e.ClarificationQuestion,
		Usage: usageResponse{
			Model:        e.Usage.Model,
			InputTokens:  e.Usage.InputTokens,
			OutputTokens: e.Usage.OutputTokens,
			CostUsd:      e.Usage.CostUsd,
		},
	}
	for _, m := range e.MatchedItems {
		resp.MatchedItems = append(resp.MatchedItems, matchedItemResponse{
			FoodId:            m.FoodId,
			FoodName:          m.FoodName,
			EstimatedQuantity: m.EstimatedQuantity,
			Unit:              m.Unit,
			Confidence:        m.Confidence,
			Assumption:        m.Assumption,
			SanityWarnings:    m.SanityWarnings,
		})
	}
	for _, s := range e.NewFoodSuggestions {
		resp.NewFoodSuggestions = append(resp.NewFoodSuggestions, newFoodSuggestionResponse{
			Name:              s.Name,
			EstimatedQuantity: s.EstimatedQuantity,
			Unit:              s.Unit,
			Confidence:        s.Confidence,
			Assumption:        s.Assumption,
			CreateParams: suggestedFoodResponse{
				MeasurementType:     s.CreateParams.MeasurementType,
				BaseUnit:            s.CreateParams.BaseUnit,
				BaseQuantity:        s.CreateParams.BaseQuantity,
				DefaultCalories:     s.CreateParams.DefaultCalories,
				DefaultProteinGrams: s.CreateParams.DefaultProteinGrams,
				DefaultCarbsGrams:   s.CreateParams.DefaultCarbsGrams,
				DefaultFatGrams:     s.CreateParams.DefaultFatGrams,
				DefaultFiberGrams:   s.CreateParams.DefaultFiberGrams,
			},
		})
	}
	return resp
}
