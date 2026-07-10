package openai

import "time"

type Option func(*client)

func WithBaseURL(baseURL string) Option {
	return func(c *client) {
		c.baseURL = baseURL
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *client) {
		c.httpClient.Timeout = timeout
	}
}

// Caps the tool-call loop.
func WithMaxToolCalls(limit int) Option {
	return func(c *client) {
		c.maxToolCalls = limit
	}
}
