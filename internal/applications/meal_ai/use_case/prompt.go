package use_case

import (
	"fmt"
	"strings"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

const systemPromptBase = `You are a nutrition assistant that estimates the foods, portions, and macros of a meal from photos and/or a text description.

Your goal is to produce a structured draft the user will review before saving. Be honest about uncertainty.

Rules:
- ALWAYS try to match each detected food to the user's existing catalog by calling the "search_foods" tool before inventing values. Prefer foods already in the catalog.
- When you match a food, you receive its stored macros. Evaluate whether those stored values are plausible for what you see; if they look wrong (e.g. grilled chicken stored at 50 kcal/100g), add a short note to that item's "sanity_warnings". Do not silently trust or silently fix them — just warn.
- If no catalog food fits (e.g. the user shows fried chicken but the catalog only has grilled chicken), do NOT force a match. Instead add a "new_food_suggestions" entry with complete "create_params" (macros per base quantity, typically per 100 g) so the user can create it.
- Estimate a quantity for every item. Use grams for solids ("g"), millilitres for liquids ("ml"), or units where natural.
- Expose your reasoning as assumptions: a short "assumption" per item, plus a global "assumptions" list (e.g. "assumed standard portion sizes", "assumed no added oil").
- Set "confidence" per item: "high" only for clearly identifiable, well-defined foods; "low" for ambiguous or mixed dishes.
- Compute "totals" as the sum across all items (matched + suggested) for the estimated quantities.
- Only set "needs_clarification" to true if the input is genuinely unusable (illegible photo, critical missing info) and put one concise question in "clarification_question"; otherwise set it to false and null respectively.`

const assumeOnlyVisibleClause = `
- Assume the user ate everything visible in the photos and ONLY what is visible. Do not invent foods that are not shown unless the user's text says so.`

func buildSystemPrompt(assumeOnlyVisible bool) string {
	if assumeOnlyVisible {
		return systemPromptBase + assumeOnlyVisibleClause
	}
	return systemPromptBase
}

// buildUserText assembles the user's instructions and any corrections from a
// prior estimate into a single message. Corrections make re-estimation a plain
// stateless call rather than a conversation.
func buildUserText(instructions string, corrections []ports.Correction) string {
	var b strings.Builder
	b.WriteString("Estimate this meal.")
	if strings.TrimSpace(instructions) != "" {
		b.WriteString("\n\nUser instructions:\n")
		b.WriteString(instructions)
	}
	if len(corrections) > 0 {
		b.WriteString("\n\nThe user corrected these assumptions from a previous estimate; honour them:")
		for _, c := range corrections {
			b.WriteString(fmt.Sprintf("\n- %s: %s", c.Item, c.Correction))
		}
	}
	return b.String()
}
