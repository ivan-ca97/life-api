package llm

import (
	"context"
	"encoding/json"
	"testing"
)

func TestCapability_Call(t *testing.T) {
	type args struct {
		Query string `json:"query"`
	}
	type result struct {
		Count int `json:"count"`
	}

	// A func literal here also verifies type inference (no explicit type args).
	handler := func(ctx context.Context, a args) (result, error) {
		return result{Count: len(a.Query)}, nil
	}
	tool := NewCapability("search", "desc", json.RawMessage(`{"type":"object"}`), handler)

	if tool.Name() != "search" || tool.Description() != "desc" {
		t.Fatalf("unexpected metadata: %q / %q", tool.Name(), tool.Description())
	}

	output, err := tool.Call(context.Background(), json.RawMessage(`{"query":"hola"}`))
	if err != nil {
		t.Fatalf("Call returned error: %v", err)
	}
	if output != `{"count":4}` {
		t.Fatalf("unexpected output: %q", output)
	}
}
