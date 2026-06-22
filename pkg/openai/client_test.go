package openai

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// roundTrip is one canned OpenAI response plus an assertion over the request
// that produced it, so we can drive the tool-call loop deterministically.
type roundTrip struct {
	assert   func(t *testing.T, req chatRequest)
	response string
}

func newMockServer(t *testing.T, trips []roundTrip) *httptest.Server {
	t.Helper()
	call := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("missing/wrong auth header: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("server could not decode request: %v", err)
		}
		if call >= len(trips) {
			t.Fatalf("unexpected extra API call #%d", call+1)
		}
		if trips[call].assert != nil {
			trips[call].assert(t, req)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, trips[call].response)
		call++
	}))
}

func newTestClient(baseURL string) *Client {
	return NewClient(Config{APIKey: "test-key", Model: "gpt-4o", BaseURL: baseURL, MaxToolCalls: 8})
}

func TestComplete_ToolLoopAndStructuredOutput(t *testing.T) {
	finalJSON := `{"totals":{"calories":330}}`

	trips := []roundTrip{
		{
			// First call: vision input + tool + response_format must all be present.
			assert: func(t *testing.T, req chatRequest) {
				if len(req.Tools) != 1 || req.Tools[0].Function.Name != "search_foods" {
					t.Errorf("expected search_foods tool, got %+v", req.Tools)
				}
				if req.ResponseFormat == nil || req.ResponseFormat.Type != "json_schema" {
					t.Errorf("expected json_schema response_format, got %+v", req.ResponseFormat)
				}
				userContent := string(req.Messages[1].Content)
				if !strings.Contains(userContent, "image_url") || !strings.Contains(userContent, "data:image/jpeg;base64,") {
					t.Errorf("expected base64 image part in user content, got %s", userContent)
				}
			},
			response: `{"choices":[{"message":{"role":"assistant","tool_calls":[{"id":"call_1","type":"function","function":{"name":"search_foods","arguments":"{\"query\":\"pollo\"}"}}]}}],"usage":{"prompt_tokens":100,"completion_tokens":20}}`,
		},
		{
			// Second call: the tool result must have been appended.
			assert: func(t *testing.T, req chatRequest) {
				last := req.Messages[len(req.Messages)-1]
				if last.Role != "tool" || last.ToolCallID != "call_1" {
					t.Errorf("expected tool result message, got %+v", last)
				}
			},
			response: `{"choices":[{"message":{"role":"assistant","content":"` + strings.ReplaceAll(finalJSON, `"`, `\"`) + `"}}],"usage":{"prompt_tokens":150,"completion_tokens":30}}`,
		},
	}

	server := newMockServer(t, trips)
	defer server.Close()

	client := newTestClient(server.URL)

	var toolCalls int
	result, err := client.Complete(context.Background(), CompletionRequest{
		System:   "you are a nutritionist",
		UserText: "estimate this meal",
		Images:   []Image{{MimeType: "image/jpeg", Data: []byte("fakebytes")}},
		Tools: []Tool{{
			Name:        "search_foods",
			Description: "search the food database",
			Parameters:  json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}},"required":["query"],"additionalProperties":false}`),
		}},
		ToolHandler: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			toolCalls++
			if name != "search_foods" {
				t.Errorf("unexpected tool %q", name)
			}
			return `[{"food_id":"abc","name":"pollo a la plancha","calories_per_100g":165}]`, nil
		},
		ResponseSchema: &ResponseSchema{
			Name:   "meal_estimate",
			Strict: true,
			Schema: json.RawMessage(`{"type":"object","properties":{"totals":{"type":"object"}},"required":["totals"],"additionalProperties":false}`),
		},
	})
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if toolCalls != 1 {
		t.Errorf("expected ToolHandler called once, got %d", toolCalls)
	}
	if result.Content != finalJSON {
		t.Errorf("expected final content %s, got %s", finalJSON, result.Content)
	}
	// Usage summed across both round trips: input 100+150, output 20+30.
	if result.Usage.InputTokens != 250 || result.Usage.OutputTokens != 50 {
		t.Errorf("unexpected usage: %+v", result.Usage)
	}
}

func TestComplete_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		io.WriteString(w, `{"error":{"message":"rate limit","type":"rate_limit_error"}}`)
	}))
	defer server.Close()

	_, err := newTestClient(server.URL).Complete(context.Background(), CompletionRequest{UserText: "hi"})
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests || apiErr.Type != "rate_limit_error" {
		t.Errorf("unexpected APIError: %+v", apiErr)
	}
}

func TestComplete_MaxToolCalls(t *testing.T) {
	// Always returns a tool call -> the loop must give up at the budget.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"choices":[{"message":{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"search_foods","arguments":"{}"}}]}}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`)
	}))
	defer server.Close()

	client := NewClient(Config{APIKey: "test-key", Model: "gpt-4o", BaseURL: server.URL, MaxToolCalls: 3})
	_, err := client.Complete(context.Background(), CompletionRequest{
		UserText:    "hi",
		Tools:       []Tool{{Name: "search_foods", Parameters: json.RawMessage(`{"type":"object"}`)}},
		ToolHandler: func(ctx context.Context, name string, args json.RawMessage) (string, error) { return "[]", nil },
	})
	if _, ok := err.(*maxToolCallsError); !ok {
		t.Fatalf("expected maxToolCallsError, got %T: %v", err, err)
	}
}

func TestCostUSD(t *testing.T) {
	got := CostUSD("gpt-4o", Usage{InputTokens: 1_000_000, OutputTokens: 1_000_000})
	if got != 12.50 {
		t.Errorf("expected 12.50, got %v", got)
	}
	if CostUSD("unknown-model", Usage{InputTokens: 1000}) != 0 {
		t.Errorf("expected 0 for unknown model")
	}
}
