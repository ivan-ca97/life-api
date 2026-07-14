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
  "required": ["queries"],
  "properties": {
    "queries": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Several short keywords (in Spanish) for ONE component food: split its name into its individual words and add a likely synonym or broader category term. Returns the deduplicated union of foods whose name contains any keyword."
    }
  }
}`)

// estimateSchema is the Structured Outputs schema for the final answer. Strict
// mode requires additionalProperties:false and every property listed in
// required; optional values are expressed as nullable types.
var estimateSchema = json.RawMessage(`{
  "type": "object",
  "additionalProperties": false,
  "required": ["matched_items", "new_food_suggestions", "totals", "assumptions", "no_food_detected", "needs_clarification", "clarification_question"],
  "properties": {
    "matched_items": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["food_id", "food_name", "estimated_quantity", "unit", "quantity_source", "confidence", "assumption", "sanity_warnings", "suggested_correction"],
        "properties": {
          "food_id": { "type": "string" },
          "food_name": { "type": "string" },
          "estimated_quantity": { "type": "number" },
          "unit": { "type": "string" },
          "quantity_source": { "type": "string", "enum": ["user", "estimated"] },
          "confidence": { "type": "string", "enum": ["low", "medium", "high"] },
          "assumption": { "type": "string" },
          "sanity_warnings": { "type": "array", "items": { "type": "string" } },
          "suggested_correction": {
            "type": ["object", "null"],
            "additionalProperties": false,
            "required": ["calories", "protein_grams", "carbs_grams", "fat_grams", "fiber_grams", "reason"],
            "properties": {
              "calories": { "type": "number" },
              "protein_grams": { "type": "number" },
              "carbs_grams": { "type": "number" },
              "fat_grams": { "type": "number" },
              "fiber_grams": { "type": "number" },
              "reason": { "type": "string" }
            }
          }
        }
      }
    },
    "new_food_suggestions": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["name", "estimated_quantity", "unit", "quantity_source", "confidence", "assumption", "create_params"],
        "properties": {
          "name": { "type": "string" },
          "estimated_quantity": { "type": "number" },
          "unit": { "type": "string" },
          "quantity_source": { "type": "string", "enum": ["user", "estimated"] },
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
    "no_food_detected": { "type": "boolean" },
    "needs_clarification": { "type": "boolean" },
    "clarification_question": { "type": ["string", "null"] }
  }
}`)

// suggestedCorrectionOutput is the model's proposed fix for a matched food whose
// stored macros look implausible (per the food's base quantity).
type suggestedCorrectionOutput struct {
	Calories     float64 `json:"calories"`
	ProteinGrams float64 `json:"protein_grams"`
	CarbsGrams   float64 `json:"carbs_grams"`
	FatGrams     float64 `json:"fat_grams"`
	FiberGrams   float64 `json:"fiber_grams"`
	Reason       string  `json:"reason"`
}

// modelOutput is the shape the model returns (guaranteed by estimateSchema).
type modelOutput struct {
	MatchedItems []struct {
		FoodId              string                     `json:"food_id"`
		FoodName            string                     `json:"food_name"`
		EstimatedQuantity   float64                    `json:"estimated_quantity"`
		Unit                string                     `json:"unit"`
		QuantitySource      string                     `json:"quantity_source"`
		Confidence          string                     `json:"confidence"`
		Assumption          string                     `json:"assumption"`
		SanityWarnings      []string                   `json:"sanity_warnings"`
		SuggestedCorrection *suggestedCorrectionOutput `json:"suggested_correction"`
	} `json:"matched_items"`
	NewFoodSuggestions []struct {
		Name              string  `json:"name"`
		EstimatedQuantity float64 `json:"estimated_quantity"`
		Unit              string  `json:"unit"`
		QuantitySource    string  `json:"quantity_source"`
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
	NoFoodDetected        bool     `json:"no_food_detected"`
	NeedsClarification    bool     `json:"needs_clarification"`
	ClarificationQuestion *string  `json:"clarification_question"`
}

func (o *modelOutput) toDomain() *domain.MealEstimate {
	estimate := &domain.MealEstimate{
		Assumptions:        o.Assumptions,
		NoFoodDetected:     o.NoFoodDetected,
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
		item := domain.MatchedItem{
			FoodId:            m.FoodId,
			FoodName:          m.FoodName,
			EstimatedQuantity: m.EstimatedQuantity,
			Unit:              m.Unit,
			QuantitySource:    m.QuantitySource,
			Confidence:        m.Confidence,
			Assumption:        m.Assumption,
			SanityWarnings:    m.SanityWarnings,
		}
		if m.SuggestedCorrection != nil {
			item.SuggestedCorrection = &domain.SuggestedCorrection{
				Calories:     m.SuggestedCorrection.Calories,
				ProteinGrams: m.SuggestedCorrection.ProteinGrams,
				CarbsGrams:   m.SuggestedCorrection.CarbsGrams,
				FatGrams:     m.SuggestedCorrection.FatGrams,
				FiberGrams:   m.SuggestedCorrection.FiberGrams,
				Reason:       m.SuggestedCorrection.Reason,
			}
		}
		estimate.MatchedItems = append(estimate.MatchedItems, item)
	}
	for _, s := range o.NewFoodSuggestions {
		item := domain.NewFoodSuggestion{
			Name:              s.Name,
			EstimatedQuantity: s.EstimatedQuantity,
			Unit:              s.Unit,
			QuantitySource:    s.QuantitySource,
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
		}
		estimate.NewFoodSuggestions = append(estimate.NewFoodSuggestions, item)
	}
	return estimate
}
