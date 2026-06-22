package openai

// ModelPricing is the USD price per 1,000,000 tokens for a model.
type ModelPricing struct {
	InputPerMTok  float64
	OutputPerMTok float64
}

// pricing holds list prices used to convert token usage into a cost estimate
// for the per-user spend limits (see docs/ai-asistente-plan.md §5).
//
// IMPORTANT: these are list prices and CHANGE over time. Verify against
// https://openai.com/api/pricing before relying on them for billing. They are
// only an estimate for soft monthly caps, not an invoice.
var pricing = map[string]ModelPricing{
	"gpt-4o":      {InputPerMTok: 2.50, OutputPerMTok: 10.00},
	"gpt-4o-mini": {InputPerMTok: 0.15, OutputPerMTok: 0.60},
	"gpt-4.1":     {InputPerMTok: 2.00, OutputPerMTok: 8.00},
}

// CostUSD estimates the cost of a Usage for the given model. Returns 0 for an
// unknown model (callers should treat 0 as "unpriced", not "free").
func CostUSD(model string, usage Usage) float64 {
	p, ok := pricing[model]
	if !ok {
		return 0
	}
	const perMillion = 1_000_000.0
	input := float64(usage.InputTokens) / perMillion * p.InputPerMTok
	output := float64(usage.OutputTokens) / perMillion * p.OutputPerMTok
	return input + output
}

// HasPricing reports whether the model has a known price entry.
func HasPricing(model string) bool {
	_, ok := pricing[model]
	return ok
}
