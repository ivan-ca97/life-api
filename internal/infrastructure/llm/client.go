package llm

import "context"

// Client is an LLM provider that completes a prompt, running the tool-call loop.
type Client interface {
	Complete(ctx context.Context, prompt Prompt) (*Result, error)
	Provider() string
	Model() string
}
