package openai

import "time"

type Option func(*client)

func WithBaseUrl(baseUrl string) Option {
	return func(c *client) {
		c.baseUrl = baseUrl
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
