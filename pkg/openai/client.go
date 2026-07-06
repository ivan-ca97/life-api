package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client is a minimal OpenAI Chat Completions client supporting vision inputs,
// tool calling, and structured outputs. It is safe for concurrent use.
type Client struct {
	config Config
	http   *http.Client
}

func NewClient(config Config) *Client {
	config = config.withDefaults()
	return &Client{
		config: config,
		http:   &http.Client{Timeout: config.Timeout},
	}
}

// Complete runs a chat completion, transparently servicing any tool calls the
// model makes via req.ToolHandler, until the model returns a final answer or
// the tool-call budget is exhausted. Usage is summed across every round trip.
func (c *Client) Complete(ctx context.Context, req CompletionRequest) (*CompletionResult, error) {
	messages := []chatMessage{
		{Role: "system", Content: mustRawString(req.System)},
		{Role: "user", Content: buildUserContent(req.UserText, req.Images)},
	}

	tools := buildTools(req.Tools)
	responseFormat := buildResponseFormat(req.ResponseSchema)

	var usage Usage
	for i := 0; i < c.config.MaxToolCalls; i++ {
		response, err := c.send(ctx, chatRequest{
			Model:          c.config.Model,
			Messages:       messages,
			Tools:          tools,
			ResponseFormat: responseFormat,
		})
		if err != nil {
			return nil, err
		}
		usage.InputTokens += response.Usage.PromptTokens
		usage.OutputTokens += response.Usage.CompletionTokens
		usage.Calls++

		if len(response.Choices) == 0 {
			return nil, fmt.Errorf("openai: response had no choices")
		}
		message := response.Choices[0].Message

		// No tool calls -> the model produced its final answer.
		if len(message.ToolCalls) == 0 {
			return &CompletionResult{Content: contentString(message.Content), Usage: usage}, nil
		}

		// Echo the assistant's tool-call message, then append one tool result
		// per call before looping back to the model.
		messages = append(messages, message)
		for _, call := range message.ToolCalls {
			if req.ToolHandler == nil {
				return nil, fmt.Errorf("openai: model requested tool %q but no ToolHandler was provided", call.Function.Name)
			}
			result, err := req.ToolHandler(ctx, call.Function.Name, json.RawMessage(call.Function.Arguments))
			if err != nil {
				return nil, fmt.Errorf("openai: tool %q failed: %w", call.Function.Name, err)
			}
			messages = append(messages, chatMessage{
				Role:       "tool",
				ToolCallID: call.ID,
				Content:    mustRawString(result),
			})
		}
	}

	return nil, &maxToolCallsError{limit: c.config.MaxToolCalls}
}

func (c *Client) send(ctx context.Context, body chatRequest) (*chatResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("openai: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("openai: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: request failed: %w", err)
	}
	defer httpResp.Body.Close()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: read response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, parseAPIError(httpResp.StatusCode, raw)
	}

	var response chatResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return nil, fmt.Errorf("openai: decode response: %w", err)
	}
	return &response, nil
}

func parseAPIError(status int, raw []byte) error {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	_ = json.Unmarshal(raw, &envelope)
	message := envelope.Error.Message
	if message == "" {
		message = string(raw)
	}
	return &APIError{StatusCode: status, Type: envelope.Error.Type, Message: message}
}

// --- request/content builders ---

func buildUserContent(text string, images []Image) json.RawMessage {
	// With no images, a plain string content is the simplest valid form.
	if len(images) == 0 {
		return mustRawString(text)
	}

	parts := make([]any, 0, len(images)+1)
	if text != "" {
		parts = append(parts, textPart{Type: "text", Text: text})
	}
	for _, img := range images {
		dataURI := fmt.Sprintf("data:%s;base64,%s", img.MimeType, base64.StdEncoding.EncodeToString(img.Data))
		part := imagePart{Type: "image_url"}
		part.ImageURL.URL = dataURI
		parts = append(parts, part)
	}
	encoded, _ := json.Marshal(parts)
	return encoded
}

func buildTools(tools []Tool) []wireTool {
	if len(tools) == 0 {
		return nil
	}
	wire := make([]wireTool, len(tools))
	for i, t := range tools {
		wire[i] = wireTool{
			Type: "function",
			Function: wireToolFunc{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		}
	}
	return wire
}

func buildResponseFormat(schema *ResponseSchema) *responseFormat {
	if schema == nil {
		return nil
	}
	return &responseFormat{
		Type: "json_schema",
		JSONSchema: &jsonSchemaSpec{
			Name:   schema.Name,
			Strict: schema.Strict,
			Schema: schema.Schema,
		},
	}
}

// mustRawString encodes s as a JSON string. Marshaling a string never fails.
func mustRawString(s string) json.RawMessage {
	encoded, _ := json.Marshal(s)
	return encoded
}

// contentString extracts a plain-string content. Final assistant messages with
// a response schema are always strings; if an array slips through we return it
// verbatim so the caller can still inspect it.
func contentString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return string(raw)
}
