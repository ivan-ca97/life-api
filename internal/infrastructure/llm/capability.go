package llm

import (
	"context"
	"encoding/json"
)

type Handler[A any, R any] func(ctx context.Context, arguments A) (R, error)

// capability adapts a typed Handler into a Tool.
type capability[A any, R any] struct {
	name        string
	description string
	schema      json.RawMessage
	call        Handler[A, R]
}

var _ Tool = (*capability[struct{}, struct{}])(nil)

func NewCapability[A any, R any](name, description string, schema json.RawMessage, call Handler[A, R]) *capability[A, R] {
	return &capability[A, R]{
		name:        name,
		description: description,
		schema:      schema,
		call:        call,
	}
}

func (c *capability[A, R]) Name() string            { return c.name }
func (c *capability[A, R]) Description() string     { return c.description }
func (c *capability[A, R]) Schema() json.RawMessage { return c.schema }

func (c *capability[A, R]) Call(ctx context.Context, arguments json.RawMessage) (string, error) {
	var args A
	err := json.Unmarshal(arguments, &args)
	if err != nil {
		return "", err
	}
	result, err := c.call(ctx, args)
	if err != nil {
		return "", err
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
