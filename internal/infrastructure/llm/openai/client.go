package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"
)

const (
	providerName = "openai"

	defaultBaseUrl      = "https://api.openai.com/v1"
	defaultTimeout      = 60 * time.Second
	defaultMaxToolCalls = 8
)

type client struct {
	apiKey       string
	model        string
	baseUrl      string
	maxToolCalls int
	httpClient   *http.Client
}

var _ llm.Client = (*client)(nil)

func NewClient(apiKey, model string, options ...Option) *client {
	c := &client{
		apiKey:       apiKey,
		model:        model,
		baseUrl:      defaultBaseUrl,
		maxToolCalls: defaultMaxToolCalls,
		httpClient:   &http.Client{Timeout: defaultTimeout},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *client) Provider() string { return providerName }
func (c *client) Model() string    { return c.model }

// Complete runs a completion, delegating the tool-call loop to an exchange.
func (c *client) Complete(ctx context.Context, prompt llm.Prompt) (*llm.Result, error) {
	ex := newExchange(c, prompt)
	result, err := ex.run(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *client) send(ctx context.Context, body chatRequest) (*chatResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("openai: marshal request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("openai: build request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("openai: request failed: %w", err)
	}
	defer httpResponse.Body.Close()

	raw, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: read response: %w", err)
	}
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		return nil, parseProviderError(httpResponse.StatusCode, raw)
	}

	var response chatResponse
	err = json.Unmarshal(raw, &response)
	if err != nil {
		return nil, fmt.Errorf("openai: decode response: %w", err)
	}
	return &response, nil
}
