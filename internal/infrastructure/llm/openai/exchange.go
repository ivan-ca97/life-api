package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"
)

type toolIndex map[string]llm.Tool

// exchange is the mutable state of one Complete call: the growing message
// history, accumulated usage, and the reusable request and tool index.
type exchange struct {
	client   *client
	request  chatRequest
	tools    toolIndex
	messages []chatMessage
	usage    llm.Usage
}

func newExchange(c *client, prompt llm.Prompt) *exchange {
	return &exchange{
		client:   c,
		request:  newChatRequest(c.model, prompt),
		tools:    indexTools(prompt.Tools),
		messages: buildMessages(prompt.Conversation),
	}
}

// run drives the tool-call loop until a final answer or the tool budget runs out.
func (e *exchange) run(ctx context.Context) (*llm.Result, error) {
	for range e.client.maxToolCalls {
		result, err := e.step(ctx)
		if err != nil {
			return nil, err
		}
		if result != nil {
			return result, nil
		}
	}

	return nil, fmt.Errorf("openai: exceeded max tool calls (%d)", e.client.maxToolCalls)
}

// step runs one round trip. It returns a non-nil result when the model gave a
// final answer, or (nil, nil) when it made tool calls and the loop should continue.
func (e *exchange) step(ctx context.Context) (*llm.Result, error) {
	e.request.Messages = e.messages
	response, err := e.client.send(ctx, e.request)
	if err != nil {
		return nil, err
	}
	e.usage.AddCall(response.Usage.PromptTokens, response.Usage.CompletionTokens)

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("openai: response had no choices")
	}

	err = finishReasonError(response.Choices[0].FinishReason)
	if err != nil {
		return nil, err
	}
	message := response.Choices[0].Message

	// The model wants to call tools: run them; the loop will continue.
	if len(message.ToolCalls) > 0 {
		results, err := e.runToolCalls(ctx, message.ToolCalls)
		if err != nil {
			return nil, err
		}
		e.messages = append(e.messages, message)
		e.messages = append(e.messages, results...)
		return nil, nil
	}

	// No tool calls: this is the final answer.
	result := llm.NewResult(contentString(message.Content), e.usage)
	return result, nil
}

func (e *exchange) runToolCalls(ctx context.Context, calls []wireToolCall) ([]chatMessage, error) {
	results := make([]chatMessage, 0, len(calls))
	for _, call := range calls {
		result, err := e.runToolCall(ctx, call)
		if err != nil {
			return nil, err
		}
		results = append(results, *result)
	}
	return results, nil
}

func (e *exchange) runToolCall(ctx context.Context, call wireToolCall) (*chatMessage, error) {
	tool, found := e.tools[call.Function.Name]
	if !found {
		return nil, fmt.Errorf("openai: model called unknown tool %q", call.Function.Name)
	}
	output, err := tool.Call(ctx, json.RawMessage(call.Function.Arguments))
	if err != nil {
		return nil, fmt.Errorf("openai: tool %q failed: %w", call.Function.Name, err)
	}
	result := chatMessage{
		Role:       "tool",
		ToolCallID: call.ID,
		Content:    mustRawString(output),
	}
	return &result, nil
}

func indexTools(tools []llm.Tool) toolIndex {
	byName := make(toolIndex, len(tools))
	for _, tool := range tools {
		byName[tool.Name()] = tool
	}
	return byName
}
