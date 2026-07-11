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
	"github.com/ivan-ca97/life/pkg/openai"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/domain"
	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

const (
	maxPhotos       = 4
	foodSearchLimit = 8

	providerOpenAI        = "openai"
	operationMealEstimate = "meal_estimate"
)

type mealEstimationUseCase struct {
	completer    ports.Completer
	foodSearch   ports.FoodSearch
	imageFetcher ports.ImageFetcher
	quota        ports.QuotaGuard
	logger       ports.InteractionLogger
	authorizer   auth.AuthorizationService
	model        string
}

var _ ports.MealEstimationUseCase = (*mealEstimationUseCase)(nil)

func NewMealEstimationUseCase(
	completer ports.Completer,
	foodSearch ports.FoodSearch,
	imageFetcher ports.ImageFetcher,
	quota ports.QuotaGuard,
	logger ports.InteractionLogger,
	authorizer auth.AuthorizationService,
	model string,
) *mealEstimationUseCase {
	return &mealEstimationUseCase{
		completer:    completer,
		foodSearch:   foodSearch,
		imageFetcher: imageFetcher,
		quota:        quota,
		logger:       logger,
		authorizer:   authorizer,
		model:        model,
	}
}

func (u *mealEstimationUseCase) Estimate(ctx context.Context, input ports.EstimateInput) (*domain.MealEstimate, error) {
	if err := u.authorizer.Authorize(ctx, input.UserId, permissions.MealsCreate); err != nil {
		return nil, err
	}
	if len(input.PhotoUrls) == 0 && strings.TrimSpace(input.Instructions) == "" {
		return nil, domain.ErrNoInput
	}
	if len(input.PhotoUrls) > maxPhotos {
		return nil, domain.ErrTooManyPhotos
	}
	if err := u.quota.CheckQuota(input.UserId); err != nil {
		return nil, err
	}

	images, err := u.fetchImages(ctx, input.PhotoUrls)
	if err != nil {
		return nil, err
	}

	request := openai.CompletionRequest{
		System:      buildSystemPrompt(input.AssumeOnlyVisible),
		UserText:    buildUserText(input.Instructions, input.Corrections),
		Images:      images,
		Tools:       []openai.Tool{{Name: searchFoodsToolName, Description: "Search the user's food catalog by name. Returns matching foods with their stored macros.", Parameters: searchFoodsParams}},
		ToolHandler: u.searchFoodsHandler(input.UserId),
		ResponseSchema: &openai.ResponseSchema{
			Name:   "meal_estimate",
			Strict: true,
			Schema: estimateSchema,
		},
	}

	start := time.Now()
	result, err := u.completer.Complete(ctx, request)
	latencyMs := int(time.Since(start).Milliseconds())
	if err != nil {
		u.logInteraction(input, providerErrorOutcome(err), openai.Usage{}, 0, latencyMs, nil)
		return nil, cerr.NewInternalError("meal ai estimation", err)
	}

	var output modelOutput
	if err := json.Unmarshal([]byte(result.Content), &output); err != nil {
		u.logInteraction(input, interactionOutcome{status: "error", errorType: "decode"}, result.Usage, 0, latencyMs, nil)
		return nil, cerr.NewInternalError("decoding meal ai estimate", err)
	}

	// Record usage against the user's monthly quota (best effort: a recording
	// failure must not discard a successful estimate).
	cost := openai.CostUSD(u.model, result.Usage)
	_ = u.quota.RecordUsage(input.UserId, ports.UsageDelta{
		Requests:     1,
		InputTokens:  int64(result.Usage.InputTokens),
		OutputTokens: int64(result.Usage.OutputTokens),
		CostUsd:      cost,
	})

	estimate := output.toDomain()
	estimate.Usage = domain.Usage{
		Model:        u.model,
		InputTokens:  int64(result.Usage.InputTokens),
		OutputTokens: int64(result.Usage.OutputTokens),
		CostUsd:      cost,
	}

	u.logInteraction(input, interactionOutcome{status: "ok"}, result.Usage, cost, latencyMs, estimate)
	return estimate, nil
}

// interactionOutcome carries the status/error-type for one logged interaction.
type interactionOutcome struct {
	status    string
	errorType string
}

// providerErrorOutcome maps a completer error to an interaction outcome, keeping
// the provider's error type (e.g. "insufficient_quota") when available.
func providerErrorOutcome(err error) interactionOutcome {
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		return interactionOutcome{status: "error", errorType: apiErr.Type}
	}
	return interactionOutcome{status: "error", errorType: "provider"}
}

// logInteraction records one interaction, best-effort (a logging failure must
// not affect the estimate). estimate is nil on failure.
func (u *mealEstimationUseCase) logInteraction(input ports.EstimateInput, outcome interactionOutcome, usage openai.Usage, cost float64, latencyMs int, estimate *domain.MealEstimate) {
	metadata := map[string]any{"photo_count": len(input.PhotoUrls)}
	if estimate != nil {
		metadata["item_count"] = len(estimate.MatchedItems)
		metadata["suggestion_count"] = len(estimate.NewFoodSuggestions)
	}
	_ = u.logger.LogInteraction(ports.InteractionEntry{
		UserId:        input.UserId,
		Operation:     operationMealEstimate,
		Provider:      providerOpenAI,
		Model:         u.model,
		Status:        outcome.status,
		ErrorType:     outcome.errorType,
		InputTokens:   int64(usage.InputTokens),
		OutputTokens:  int64(usage.OutputTokens),
		CostUsd:       cost,
		LatencyMs:     latencyMs,
		ProviderCalls: usage.Calls,
		InputSummary:  input.Instructions,
		Metadata:      metadata,
	})
}

func (u *mealEstimationUseCase) fetchImages(ctx context.Context, urls []string) ([]openai.Image, error) {
	images := make([]openai.Image, 0, len(urls))
	for _, url := range urls {
		img, err := u.imageFetcher.Fetch(ctx, url)
		if err != nil {
			return nil, cerr.NewBadRequestError("could not load one of the photos")
		}
		images = append(images, openai.Image{MimeType: img.MimeType, Data: img.Data})
	}
	return images, nil
}

// searchFoodsHandler returns the tool callback bound to a user, translating the
// model's search request into a catalog lookup and back into JSON candidates.
func (u *mealEstimationUseCase) searchFoodsHandler(userId uuid.UUID) openai.ToolHandler {
	return func(ctx context.Context, name string, arguments json.RawMessage) (string, error) {
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(arguments, &args); err != nil {
			return "", err
		}
		candidates, err := u.foodSearch.Search(userId, args.Query, foodSearchLimit)
		if err != nil {
			return "", err
		}
		return marshalCandidates(candidates), nil
	}
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

func marshalCandidates(candidates []ports.FoodCandidate) string {
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
	encoded, err := json.Marshal(out)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}
