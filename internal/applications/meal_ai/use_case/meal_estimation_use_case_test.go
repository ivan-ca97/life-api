package use_case

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"

	"github.com/ivan-ca97/life/internal/applications/meal_ai/ports"
)

const mockCost = 0.5

type mockClient struct {
	toolCalled bool
}

func (m *mockClient) Provider() string { return "openai" }
func (m *mockClient) Model() string    { return "gpt-4o" }

func (m *mockClient) Complete(ctx context.Context, prompt llm.Prompt) (*llm.Result, error) {
	// Simulate the model calling search_foods once, then returning the final draft.
	if len(prompt.Tools) > 0 {
		_, err := prompt.Tools[0].Call(ctx, json.RawMessage(`{"query":"chicken"}`))
		if err != nil {
			return nil, err
		}
		m.toolCalled = true
	}
	final := `{
	  "matched_items": [{"food_id":"abc","food_name":"Grilled chicken","estimated_quantity":200,"unit":"g","confidence":"medium","assumption":"no oil","sanity_warnings":[]}],
	  "new_food_suggestions": [],
	  "totals": {"calories":330,"protein_grams":62,"carbs_grams":0,"fat_grams":7,"fiber_grams":0},
	  "assumptions": ["standard portions"],
	  "needs_clarification": false,
	  "clarification_question": null
	}`
	return &llm.Result{Content: final, Usage: llm.Usage{InputTokens: 100, OutputTokens: 50, Calls: 2}}, nil
}

type mockFoodSearch struct{ called bool }

func (m *mockFoodSearch) Search(userId uuid.UUID, query string, limit int) ([]ports.FoodCandidate, error) {
	m.called = true
	cal := 165.0
	return []ports.FoodCandidate{{Id: uuid.New(), Name: "Grilled chicken", DefaultCalories: &cal, MeasurementType: "mass", BaseQuantity: 100, BaseUnit: "g"}}, nil
}

type mockImageFetcher struct{}

func (m *mockImageFetcher) Fetch(ctx context.Context, url string) (llm.Image, error) {
	return llm.Image{MediaType: "image/jpeg", Data: []byte("bytes")}, nil
}

type mockQuota struct {
	checked  bool
	recorded *ports.UsageDelta
}

func (m *mockQuota) CheckQuota(userId uuid.UUID) error { m.checked = true; return nil }
func (m *mockQuota) RecordUsage(userId uuid.UUID, delta ports.UsageDelta) error {
	m.recorded = &delta
	return nil
}

type mockPricer struct{}

func (mockPricer) CostUSD(provider, model string, inputTokens, outputTokens int64, at time.Time) (float64, error) {
	return mockCost, nil
}

type stubAuthorizer struct{}

func (stubAuthorizer) Authorize(ctx context.Context, ownerId uuid.UUID, permission string) error {
	return nil
}
func (stubAuthorizer) AuthorizeAdmin(ctx context.Context) error { return nil }

type mockLogger struct{ entry *ports.InteractionEntry }

func (m *mockLogger) LogInteraction(entry ports.InteractionEntry) error {
	m.entry = &entry
	return nil
}

func TestEstimate_HappyPath(t *testing.T) {
	client := &mockClient{}
	foodSearch := &mockFoodSearch{}
	quota := &mockQuota{}
	logger := &mockLogger{}

	uc := NewMealEstimationUseCase(client, foodSearch, &mockImageFetcher{}, quota, logger, mockPricer{}, stubAuthorizer{})

	estimate, err := uc.Estimate(context.Background(), ports.EstimateInput{
		UserId:            uuid.New(),
		PhotoUrls:         []string{"https://example.com/meal.jpg"},
		AssumeOnlyVisible: true,
	})
	if err != nil {
		t.Fatalf("Estimate returned error: %v", err)
	}

	if !quota.checked {
		t.Error("quota was not checked before spending")
	}
	if !client.toolCalled || !foodSearch.called {
		t.Error("search_foods tool was not exercised")
	}
	if len(estimate.MatchedItems) != 1 || estimate.MatchedItems[0].FoodId != "abc" {
		t.Fatalf("unexpected matched items: %+v", estimate.MatchedItems)
	}
	if estimate.Totals.Calories != 330 {
		t.Errorf("expected totals.calories 330, got %v", estimate.Totals.Calories)
	}
	if estimate.Usage.CostUsd != mockCost {
		t.Errorf("expected cost %v, got %v", mockCost, estimate.Usage.CostUsd)
	}
	if quota.recorded == nil || quota.recorded.CostUsd != mockCost {
		t.Errorf("usage not recorded correctly: %+v", quota.recorded)
	}

	// The interaction must be logged with the right metadata.
	if logger.entry == nil {
		t.Fatal("interaction was not logged")
	}
	e := logger.entry
	if e.Operation != "meal_estimate" || e.Provider != "openai" || e.Model != "gpt-4o" {
		t.Errorf("unexpected interaction identity: %+v", e)
	}
	if e.Status != "ok" || e.CostUsd != mockCost || e.ProviderCalls != 2 {
		t.Errorf("unexpected interaction metrics: status=%s cost=%v calls=%d", e.Status, e.CostUsd, e.ProviderCalls)
	}
	if e.Metadata["photo_count"] != 1 || e.Metadata["item_count"] != 1 || e.Metadata["suggestion_count"] != 0 {
		t.Errorf("unexpected interaction metadata: %+v", e.Metadata)
	}
}

func TestEstimate_NoInput(t *testing.T) {
	uc := NewMealEstimationUseCase(&mockClient{}, &mockFoodSearch{}, &mockImageFetcher{}, &mockQuota{}, &mockLogger{}, mockPricer{}, stubAuthorizer{})
	_, err := uc.Estimate(context.Background(), ports.EstimateInput{UserId: uuid.New()})
	if err == nil {
		t.Fatal("expected an error when no photos and no instructions are given")
	}
}
