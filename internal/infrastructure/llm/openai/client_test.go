package openai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"
)

type roundTrip struct {
	assert   func(t *testing.T, request chatRequest)
	response string
}

func newMockServer(t *testing.T, trips []roundTrip) *httptest.Server {
	t.Helper()
	call := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("missing/wrong auth header: %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var request chatRequest
		err := json.Unmarshal(body, &request)
		if err != nil {
			t.Fatalf("server could not decode request: %v", err)
		}
		if call >= len(trips) {
			t.Fatalf("unexpected extra API call #%d", call+1)
		}
		if trips[call].assert != nil {
			trips[call].assert(t, request)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, trips[call].response)
		call++
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestComplete_ToolLoopAndStructuredOutput(t *testing.T) {
	finalContent := `{"ok":true}`

	trips := []roundTrip{
		{
			assert: func(t *testing.T, request chatRequest) {
				if len(request.Tools) != 1 || request.Tools[0].Function.Name != "search_foods" {
					t.Errorf("expected search_foods tool, got %+v", request.Tools)
				}
				if request.ResponseFormat == nil || request.ResponseFormat.Type != "json_schema" {
					t.Errorf("expected json_schema response_format, got %+v", request.ResponseFormat)
				}
				userContent := string(request.Messages[1].Content)
				if !strings.Contains(userContent, "image_url") || !strings.Contains(userContent, "data:image/jpeg;base64,") {
					t.Errorf("expected base64 image part in user content, got %s", userContent)
				}
			},
			response: `{"choices":[{"message":{"role":"assistant","tool_calls":[{"id":"call_1","type":"function","function":{"name":"search_foods","arguments":"{\"query\":\"chicken\"}"}}]}}],"usage":{"prompt_tokens":100,"completion_tokens":20}}`,
		},
		{
			assert: func(t *testing.T, request chatRequest) {
				last := request.Messages[len(request.Messages)-1]
				if last.Role != "tool" || last.ToolCallID != "call_1" {
					t.Errorf("expected tool result message, got %+v", last)
				}
			},
			response: `{"choices":[{"message":{"role":"assistant","content":"` + strings.ReplaceAll(finalContent, `"`, `\"`) + `"}}],"usage":{"prompt_tokens":150,"completion_tokens":30}}`,
		},
	}

	server := newMockServer(t, trips)
	defer server.Close()

	toolCalls := 0
	type searchArgs struct {
		Query string `json:"query"`
	}
	searchTool := llm.NewCapability("search_foods", "search the food catalog",
		json.RawMessage(`{"type":"object","additionalProperties":false,"required":["query"],"properties":{"query":{"type":"string"}}}`),
		func(ctx context.Context, args searchArgs) ([]map[string]any, error) {
			toolCalls++
			if args.Query == "" {
				t.Errorf("expected a non-empty query")
			}
			return []map[string]any{{"food_id": "abc"}}, nil
		})

	prompt := llm.Prompt{
		Conversation: llm.SingleTurn("you are a nutritionist", "estimate this meal", []llm.Image{{MediaType: "image/jpeg", Data: []byte("fakebytes")}}),
		Tools:        []llm.Tool{searchTool},
		ResponseSchema: &llm.ResponseSchema{
			Name:   "meal_estimate",
			Strict: true,
			Schema: json.RawMessage(`{"type":"object"}`),
		},
	}

	client := NewClient("test-key", "gpt-4o", WithBaseURL(server.URL), WithMaxToolCalls(8))
	result, err := client.Complete(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}

	if toolCalls != 1 {
		t.Errorf("expected the tool called once, got %d", toolCalls)
	}
	if result.Content != finalContent {
		t.Errorf("expected content %s, got %s", finalContent, result.Content)
	}
	if result.Usage.InputTokens != 250 || result.Usage.OutputTokens != 50 || result.Usage.Calls != 2 {
		t.Errorf("unexpected usage: %+v", result.Usage)
	}
}

func TestComplete_ProviderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		io.WriteString(w, `{"error":{"message":"quota","type":"insufficient_quota"}}`)
	}))
	defer server.Close()

	client := NewClient("test-key", "gpt-4o", WithBaseURL(server.URL))
	_, err := client.Complete(context.Background(), llm.Prompt{Conversation: llm.SingleTurn("", "hi", nil)})

	var providerError *llm.ProviderError
	if !errors.As(err, &providerError) {
		t.Fatalf("expected *llm.ProviderError, got %T: %v", err, err)
	}
	if providerError.Provider != "openai" || providerError.StatusCode != http.StatusTooManyRequests || providerError.Type != "insufficient_quota" {
		t.Errorf("unexpected provider error: %+v", providerError)
	}
}

func TestComplete_Truncated(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"choices":[{"message":{"role":"assistant","content":"{\"partial\":"},"finish_reason":"length"}],"usage":{"prompt_tokens":10,"completion_tokens":5}}`)
	}))
	defer server.Close()

	client := NewClient("test-key", "gpt-4o", WithBaseURL(server.URL))
	_, err := client.Complete(context.Background(), llm.Prompt{Conversation: llm.SingleTurn("", "hi", nil)})
	if !errors.Is(err, llm.ErrTruncated) {
		t.Fatalf("expected llm.ErrTruncated, got %v", err)
	}
}

func TestComplete_MaxToolCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"choices":[{"message":{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"loop","arguments":"{}"}}]}}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`)
	}))
	defer server.Close()

	loopTool := llm.NewCapability("loop", "loops", json.RawMessage(`{"type":"object"}`),
		func(ctx context.Context, args struct{}) (string, error) { return "again", nil })

	client := NewClient("test-key", "gpt-4o", WithBaseURL(server.URL), WithMaxToolCalls(3))
	prompt := llm.Prompt{
		Conversation: llm.SingleTurn("", "hi", nil),
		Tools:        []llm.Tool{loopTool},
	}
	_, err := client.Complete(context.Background(), prompt)
	if err == nil || !strings.Contains(err.Error(), "max tool calls") {
		t.Fatalf("expected a max tool calls error, got %v", err)
	}
}
