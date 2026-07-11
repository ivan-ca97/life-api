package use_case

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/domain"
	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

const (
	maxPhotos       = 4
	foodSearchLimit = 8

	operationMealEstimate = "meal_estimate"

	searchFoodsDescription = "Search the user's food catalog by name. Returns matching foods with their stored macros."
)

type mealEstimationUseCase struct {
	client       llm.Client
	foodSearch   ports.FoodSearch
	imageFetcher ports.ImageFetcher
	quota        ports.QuotaGuard
	logger       ports.InteractionLogger
	pricer       ports.Pricer
	authorizer   auth.AuthorizationService
}

var _ ports.MealEstimationUseCase = (*mealEstimationUseCase)(nil)

func NewMealEstimationUseCase(
	client llm.Client,
	foodSearch ports.FoodSearch,
	imageFetcher ports.ImageFetcher,
	quota ports.QuotaGuard,
	logger ports.InteractionLogger,
	pricer ports.Pricer,
	authorizer auth.AuthorizationService,
) *mealEstimationUseCase {
	return &mealEstimationUseCase{
		client:       client,
		foodSearch:   foodSearch,
		imageFetcher: imageFetcher,
		quota:        quota,
		logger:       logger,
		pricer:       pricer,
		authorizer:   authorizer,
	}
}

func (uc *mealEstimationUseCase) Estimate(ctx context.Context, input ports.EstimateInput) (*domain.MealEstimate, error) {
	if err := uc.authorizer.Authorize(ctx, input.UserId, permissions.MealsCreate); err != nil {
		return nil, err
	}
	if len(input.PhotoUrls) == 0 && strings.TrimSpace(input.Instructions) == "" {
		return nil, domain.ErrNoInput
	}
	if len(input.PhotoUrls) > maxPhotos {
		return nil, domain.ErrTooManyPhotos
	}
	if err := uc.quota.CheckQuota(input.UserId); err != nil {
		return nil, err
	}

	images, err := uc.fetchImages(ctx, input.PhotoUrls)
	if err != nil {
		return nil, err
	}

	responseSchema := &llm.ResponseSchema{
		Name:   "meal_estimate",
		Strict: true,
		Schema: estimateSchema,
	}
	prompt := llm.Prompt{
		Conversation:   llm.SingleTurn(buildSystemPrompt(input.AssumeOnlyVisible), buildUserText(input.Instructions, input.Corrections), images),
		Tools:          []llm.Tool{uc.searchFoodsTool(input.UserId)},
		ResponseSchema: responseSchema,
	}

	start := time.Now()
	result, err := uc.client.Complete(ctx, prompt)
	latencyMs := int(time.Since(start).Milliseconds())
	if err != nil {
		uc.logInteraction(input, providerErrorOutcome(err), llm.Usage{}, 0, latencyMs, nil)
		return nil, cerr.NewInternalError("meal ai estimation", err)
	}

	var output modelOutput
	if err := json.Unmarshal([]byte(result.Content), &output); err != nil {
		uc.logInteraction(input, errorOutcome("decode"), result.Usage, 0, latencyMs, nil)
		return nil, cerr.NewInternalError("decoding meal ai estimate", err)
	}

	cost := uc.cost(result.Usage)
	// Record usage against the user's monthly quota (best effort: a recording
	// failure must not discard a successful estimate).
	delta := ports.UsageDelta{
		Requests:     1,
		InputTokens:  int64(result.Usage.InputTokens),
		OutputTokens: int64(result.Usage.OutputTokens),
		CostUsd:      cost,
	}
	_ = uc.quota.RecordUsage(input.UserId, delta)

	estimate := output.toDomain()
	estimate.Usage = domain.Usage{
		Model:        uc.client.Model(),
		InputTokens:  int64(result.Usage.InputTokens),
		OutputTokens: int64(result.Usage.OutputTokens),
		CostUsd:      cost,
	}

	uc.logInteraction(input, okOutcome(), result.Usage, cost, latencyMs, estimate)
	return estimate, nil
}

// cost prices the usage best-effort; a pricing failure yields 0.
func (uc *mealEstimationUseCase) cost(usage llm.Usage) float64 {
	cost, err := uc.pricer.CostUsd(uc.client.Provider(), uc.client.Model(), int64(usage.InputTokens), int64(usage.OutputTokens), time.Now())
	if err != nil {
		return 0
	}
	return cost
}

// interactionOutcome carries the status/error-type for one logged interaction.
type interactionOutcome struct {
	status    string
	errorType string
}

func okOutcome() interactionOutcome {
	return interactionOutcome{
		status: "ok",
	}
}

func errorOutcome(errorType string) interactionOutcome {
	return interactionOutcome{
		status:    "error",
		errorType: errorType,
	}
}

// providerErrorOutcome maps a completer error to an interaction outcome, keeping
// the provider's error type (e.g. "insufficient_quota") when available.
func providerErrorOutcome(err error) interactionOutcome {
	var providerErr *llm.ProviderError
	if errors.As(err, &providerErr) {
		return errorOutcome(providerErr.Type)
	}
	return errorOutcome("provider")
}

// logInteraction records one interaction, best-effort (a logging failure must
// not affect the estimate). estimate is nil on failure.
func (uc *mealEstimationUseCase) logInteraction(input ports.EstimateInput, outcome interactionOutcome, usage llm.Usage, cost float64, latencyMs int, estimate *domain.MealEstimate) {
	metadata := map[string]any{"photo_count": len(input.PhotoUrls)}
	if estimate != nil {
		metadata["item_count"] = len(estimate.MatchedItems)
		metadata["suggestion_count"] = len(estimate.NewFoodSuggestions)
	}
	entry := ports.InteractionEntry{
		UserId:        input.UserId,
		Operation:     operationMealEstimate,
		Provider:      uc.client.Provider(),
		Model:         uc.client.Model(),
		Status:        outcome.status,
		ErrorType:     outcome.errorType,
		InputTokens:   int64(usage.InputTokens),
		OutputTokens:  int64(usage.OutputTokens),
		CostUsd:       cost,
		LatencyMs:     latencyMs,
		ProviderCalls: usage.Calls,
		InputSummary:  input.Instructions,
		Metadata:      metadata,
	}
	_ = uc.logger.LogInteraction(entry)
}

func (uc *mealEstimationUseCase) fetchImages(ctx context.Context, urls []string) ([]llm.Image, error) {
	images := make([]llm.Image, 0, len(urls))
	for _, url := range urls {
		image, err := uc.imageFetcher.Fetch(ctx, url)
		if err != nil {
			return nil, cerr.NewBadRequestError("could not load one of the photos")
		}
		images = append(images, image)
	}
	return images, nil
}

type searchFoodsArgs struct {
	Query string `json:"query"`
}

// searchFoodsTool is the capability the model calls to look up the user's catalog.
func (uc *mealEstimationUseCase) searchFoodsTool(userId uuid.UUID) llm.Tool {
	call := func(ctx context.Context, args searchFoodsArgs) ([]candidateJSON, error) {
		candidates, err := uc.foodSearch.Search(userId, args.Query, foodSearchLimit)
		if err != nil {
			return nil, err
		}
		return toCandidateJSON(candidates), nil
	}
	return llm.NewCapability(searchFoodsToolName, searchFoodsDescription, searchFoodsParams, call)
}

type candidateJSON struct {
	FoodId          string   `json:"food_id"`
	Name            string   `json:"name"`
	Calories        *float64 `json:"default_calories"`
	ProteinGrams    *float64 `json:"default_protein_grams"`
	CarbsGrams      *float64 `json:"default_carbs_grams"`
	FatGrams        *float64 `json:"default_fat_grams"`
	FiberGrams      *float64 `json:"default_fiber_grams"`
	MeasurementType string   `json:"measurement_type"`
	BaseQuantity    float64  `json:"base_quantity"`
	BaseUnit        string   `json:"base_unit"`
}

func toCandidateJSON(candidates []ports.FoodCandidate) []candidateJSON {
	out := make([]candidateJSON, len(candidates))
	for i, c := range candidates {
		out[i] = candidateJSON{
			FoodId:          c.Id.String(),
			Name:            c.Name,
			Calories:        c.DefaultCalories,
			ProteinGrams:    c.DefaultProteinGrams,
			CarbsGrams:      c.DefaultCarbsGrams,
			FatGrams:        c.DefaultFatGrams,
			FiberGrams:      c.DefaultFiberGrams,
			MeasurementType: c.MeasurementType,
			BaseQuantity:    c.BaseQuantity,
			BaseUnit:        c.BaseUnit,
		}
	}
	return out
}
