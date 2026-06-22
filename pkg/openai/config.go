package openai

import "time"

// Config holds the settings needed to talk to the OpenAI API. It is built from
// environment variables in cmd/server/main.go and passed to NewClient.
type Config struct {
	APIKey string
	// Model is the chat-completions model id, e.g. "gpt-4o".
	Model string
	// Timeout bounds a single Complete call (including the whole tool-call loop).
	Timeout time.Duration
	// BaseURL defaults to https://api.openai.com/v1 when empty.
	BaseURL string
	// MaxToolCalls caps the number of tool-call round trips per Complete call,
	// a safety net against runaway loops and cost. Defaults to 8 when <= 0.
	MaxToolCalls int
}

const (
	defaultBaseURL      = "https://api.openai.com/v1"
	defaultTimeout      = 60 * time.Second
	defaultMaxToolCalls = 8
)

func (c Config) withDefaults() Config {
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.Timeout <= 0 {
		c.Timeout = defaultTimeout
	}
	if c.MaxToolCalls <= 0 {
		c.MaxToolCalls = defaultMaxToolCalls
	}
	return c
}

// Enabled reports whether the client has enough configuration to be used. The
// server uses this to degrade gracefully: with no API key the AI feature is
// simply not registered and the rest of the app boots normally.
func (c Config) Enabled() bool {
	return c.APIKey != "" && c.Model != ""
}
