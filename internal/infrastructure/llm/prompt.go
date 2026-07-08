package llm

import "encoding/json"

type Prompt struct {
	Conversation   Conversation
	Tools          []Tool
	ResponseSchema *ResponseSchema
}

// ResponseSchema forces the final answer to match a JSON Schema. Strict requires
// additionalProperties false and every property listed as required.
type ResponseSchema struct {
	Name   string
	Schema json.RawMessage
	Strict bool
}
