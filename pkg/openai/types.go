package openai

import (
	"context"
	"encoding/json"
	"fmt"
)

// Image is a single image input for a vision request. Data is the raw bytes;
// the client base64-encodes it into a data URI so the bucket never needs to be
// public.
type Image struct {
	MimeType string
	Data     []byte
}

// Tool describes a function the model may call during a Complete call.
// Parameters is a JSON Schema object describing the function arguments.
type Tool struct {
	Name        string
	Description string
	Parameters  json.RawMessage
}

// ToolHandler executes a tool the model asked for and returns the result as a
// string (typically JSON) that is fed back to the model. The use case provides
// this; pkg/openai never touches the database directly.
type ToolHandler func(ctx context.Context, name string, arguments json.RawMessage) (string, error)

// ResponseSchema forces the model's final answer to match a JSON Schema via
// OpenAI Structured Outputs, so the caller can json.Unmarshal it directly with
// no text parsing. Strict requires additionalProperties:false and every
// property listed in "required".
type ResponseSchema struct {
	Name   string
	Strict bool
	Schema json.RawMessage
}

// CompletionRequest is the high-level input to Complete.
type CompletionRequest struct {
	System         string
	UserText       string
	Images         []Image
	Tools          []Tool
	ToolHandler    ToolHandler
	ResponseSchema *ResponseSchema
}

// Usage is the token accounting for a whole Complete call, summed across every
// round trip in the tool-call loop.
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// CompletionResult is the output of Complete. Content is the final assistant
// message (a JSON document when ResponseSchema was set).
type CompletionResult struct {
	Content string
	Usage   Usage
}

// APIError is a non-2xx response from OpenAI.
type APIError struct {
	StatusCode int
	Type       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("openai api error (status %d, type %q): %s", e.StatusCode, e.Type, e.Message)
}

// ErrMaxToolCalls is returned when the model keeps requesting tool calls past
// Config.MaxToolCalls without producing a final answer.
type maxToolCallsError struct{ limit int }

func (e *maxToolCallsError) Error() string {
	return fmt.Sprintf("openai: exceeded max tool calls (%d) without a final answer", e.limit)
}
