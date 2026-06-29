package use_case

import (
	"encoding/json"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/domain"
)

// searchFoodsTool is the function the model calls to look up the user's catalog.
const searchFoodsToolName = "search_foods"

var searchFoodsParams = json.RawMessage(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["query"],
  "properties": {
    "query": {
      "type": "string",
      "description": "A food name to search for in the user's catalog, e.g. 'grilled chicken' or 'white rice'."
    }
  }
}`)

// estimateSchema is the Structured Outputs schema for the final answer. Strict
// mode requires additionalProperties:false and every property listed in
// required; optional values are expressed as nullable types.
var estimateSchema = json.RawMessage(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["matched_items", "new_food_suggestions", "totals", "assumptions", "needs_clarification", "clarification_question"],
  "properties": {
    "matched_items": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["food_id", "food_name", "estimated_quantity", "unit", "confidence", "assumption", "sanity_warnings"],
        "properties": {
          "food_id": { "type": "string" },
          "food_name": { "type": "string" },
          "estimated_quantity": { "type": "number" },
          "unit": { "type": "string" },
          "confidence": { "type": "string", "enum": ["low", "medium", "high"] },
          "assumption": { "type": "string" },
          "sanity_warnings": { "type": "array", "items": { "type": "string" } }
        }
      }
    },
    "new_food_suggestions": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["name", "estimated_quantity", "unit", "confidence", "assumption", "create_params"],
        "properties": {
          "name": { "type": "string" },
          "estimated_quantity": { "type": "number" },
          "unit": { "type": "string" },
          "confidence": { "type": "string", "enum": ["low", "medium", "high"] },
          "assumption": { "type": "string" },
          "create_params": {
            "type": "object",
            "additionalProperties": false,
            "required": ["measurement_type", "base_unit", "base_quantity", "default_calories", "default_protein_grams", "default_carbs_grams", "default_fat_grams", "default_fiber_grams"],
            "properties": {
              "measurement_type": { "type": "string", "enum": ["mass", "volume", "unit"] },
              "base_unit": { "type": "string" },
              "base_quantity": { "type": "number" },
              "default_calories": { "type": ["number", "null"] },
              "default_protein_grams": { "type": ["number", "null"] },
              "default_carbs_grams": { "type": ["number", "null"] },
              "default_fat_grams": { "type": ["number", "null"] },
              "default_fiber_grams": { "type": ["number", "null"] }
            }
          }
        }
      }
    },
    "totals": {
      "type": "object",
      "additionalProperties": false,
      "required": ["calories", "protein_grams", "carbs_grams", "fat_grams", "fiber_grams"],
      "properties": {
        "calories": { "type": "number" },
        "protein_grams": { "type": "number" },
        "carbs_grams": { "type": "number" },
        "fat_grams": { "type": "number" },
        "fiber_grams": { "type": "number" }
      }
    },
    "assumptions": { "type": "array", "items": { "type": "string" } },
    "needs_clarification": { "type": "boolean" },
    "clarification_question": { "type": ["string", "null"] }
  }
}`)

// modelOutput is the shape the model returns (guaranteed by estimateSchema).
type modelOutput struct {
	MatchedItems []struct {
		FoodId            string   `json:"food_id"`
		FoodName          string   `json:"food_name"`
		EstimatedQuantity float64  `json:"estimated_quantity"`
		Unit              string   `json:"unit"`
		Confidence        string   `json:"confidence"`
		Assumption        string   `json:"assumption"`
		SanityWarnings    []string `json:"sanity_warnings"`
	} `json:"matched_items"`
	NewFoodSuggestions []struct {
		Name              string  `json:"name"`
		EstimatedQuantity float64 `json:"estimated_quantity"`
		Unit              string  `json:"unit"`
		Confidence        string  `json:"confidence"`
		Assumption        string  `json:"assumption"`
		CreateParams      struct {
			MeasurementType     string   `json:"measurement_type"`
			BaseUnit            string   `json:"base_unit"`
			BaseQuantity        float64  `json:"base_quantity"`
			DefaultCalories     *float64 `json:"default_calories"`
			DefaultProteinGrams *float64 `json:"default_protein_grams"`
			DefaultCarbsGrams   *float64 `json:"default_carbs_grams"`
			DefaultFatGrams     *float64 `json:"default_fat_grams"`
			DefaultFiberGrams   *float64 `json:"default_fiber_grams"`
		} `json:"create_params"`
	} `json:"new_food_suggestions"`
	Totals struct {
		Calories     float64 `json:"calories"`
		ProteinGrams float64 `json:"protein_grams"`
		CarbsGrams   float64 `json:"carbs_grams"`
		FatGrams     float64 `json:"fat_grams"`
		FiberGrams   float64 `json:"fiber_grams"`
	} `json:"totals"`
	Assumptions           []string `json:"assumptions"`
	NeedsClarification    bool     `json:"needs_clarification"`
	ClarificationQuestion *string  `json:"clarification_question"`
}

func (o *modelOutput) toDomain() *domain.MealEstimate {
	estimate := &domain.MealEstimate{
		Assumptions:        o.Assumptions,
		NeedsClarification: o.NeedsClarification,
		Totals: domain.Totals{
			Calories:     o.Totals.Calories,
			ProteinGrams: o.Totals.ProteinGrams,
			CarbsGrams:   o.Totals.CarbsGrams,
			FatGrams:     o.Totals.FatGrams,
			FiberGrams:   o.Totals.FiberGrams,
		},
	}
	if o.ClarificationQuestion != nil {
		estimate.ClarificationQuestion = *o.ClarificationQuestion
	}
	for _, m := range o.MatchedItems {
		estimate.MatchedItems = append(estimate.MatchedItems, domain.MatchedItem{
			FoodId:            m.FoodId,
			FoodName:          m.FoodName,
			EstimatedQuantity: m.EstimatedQuantity,
			Unit:              m.Unit,
			Confidence:        m.Confidence,
			Assumption:        m.Assumption,
			SanityWarnings:    m.SanityWarnings,
		})
	}
	for _, s := range o.NewFoodSuggestions {
		estimate.NewFoodSuggestions = append(estimate.NewFoodSuggestions, domain.NewFoodSuggestion{
			Name:              s.Name,
			EstimatedQuantity: s.EstimatedQuantity,
			Unit:              s.Unit,
			Confidence:        s.Confidence,
			Assumption:        s.Assumption,
			CreateParams: domain.SuggestedFood{
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
	return estimate
}
