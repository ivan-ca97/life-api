package use_case

import (
	"fmt"
	"strings"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

const systemPromptBase = `You are a nutrition assistant that estimates the foods, portions, and macros of a meal from photos and/or a text description.

Your goal is to produce a structured draft the user will review before saving. Be honest about uncertainty.

Rules:
- Decompose the meal into its individual component foods, even when they are plated together or mixed: pasta with a meat-and-tomato sauce and grated cheese is several items (e.g. "Fideos" + the meat + "Queso rallado"), NOT one combined dish. Component foods are eaten in many other combinations, so keeping them atomic makes them reusable and matchable to the catalog. Only treat something as a single item when it is cooked as an indivisible preparation (e.g. empanada, stew, pie, breaded cutlet/milanesa).
- ALWAYS try to match each component food to the user's existing catalog by calling the "search_foods" tool before inventing values. Search with a SHORT, GENERIC term for the component ("lomo", "fideos", "queso rallado", "milanesa") — never a long composite description, which matches nothing. Prefer an existing catalog match over a new suggestion.
- When you match a food, you receive its stored macros. Evaluate whether those stored values are plausible for what you see; if they look wrong (e.g. grilled chicken stored at 50 kcal/100g), add a short note to that item's "sanity_warnings". Do not silently trust or silently fix them — just warn.
- If no catalog food fits (e.g. the user shows fried chicken but the catalog only has grilled chicken), do NOT force a match. Instead add a "new_food_suggestions" entry with complete "create_params" (macros per base quantity, typically per 100 g) so the user can create it.
- Estimate a quantity for every item. Use grams for solids ("g"), millilitres for liquids ("ml"), or units where natural.
- Expose your reasoning as assumptions: a short "assumption" per item, plus a global "assumptions" list (e.g. "assumed standard portion sizes", "assumed no added oil").
- Set "confidence" per item: "high" only for clearly identifiable, well-defined foods; "low" for ambiguous or mixed dishes.
- Compute "totals" as the sum across all items (matched + suggested) for the estimated quantities.
- Keep calories consistent with macros: for every item, and for "totals", calories must approximately equal 4*protein_grams + 4*carbs_grams + 9*fat_grams. When you fill "create_params" for a new food, set its "default_calories" to match its default macros the same way.
- Set "no_food_detected" to true if the provided photos do not appear to show any food or drink (e.g. a screenshot, a receipt, a person, an unrelated object). You may still estimate from the user's text description if it names foods; only when there is no usable food information at all should you also set "needs_clarification". When the photos clearly show food, set it to false.
- Only set "needs_clarification" to true if the input is genuinely unusable (illegible photo, no food anywhere, critical missing info) and put one concise question in "clarification_question"; otherwise set it to false and null respectively.`

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
