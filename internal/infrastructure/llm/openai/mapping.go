package openai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ivan-ca97/life/internal/infrastructure/llm"
)

func buildMessages(conversation llm.Conversation) []chatMessage {
	messages := make([]chatMessage, 0, len(conversation.Messages)+1)
	if conversation.Instructions != "" {
		system := chatMessage{
			Role:    "system",
			Content: mustRawString(conversation.Instructions),
		}
		messages = append(messages, system)
	}
	for _, message := range conversation.Messages {
		mapped := chatMessage{
			Role:    openAIRole(message.Role),
			Content: messageContent(message.Text, message.Images),
		}
		messages = append(messages, mapped)
	}
	return messages
}

func openAIRole(role llm.Role) string {
	switch role {
	case llm.RoleAssistant:
		return "assistant"
	case llm.RoleTool:
		return "tool"
	default:
		return "user"
	}
}

// messageContent is a plain string with no images, or an array of typed parts
// (text + base64 images) otherwise.
func messageContent(text string, images []llm.Image) json.RawMessage {
	if len(images) == 0 {
		return mustRawString(text)
	}
	parts := make([]any, 0, len(images)+1)
	if text != "" {
		parts = append(parts, textPart{Type: "text", Text: text})
	}
	for _, image := range images {
		dataUri := fmt.Sprintf("data:%s;base64,%s", image.MediaType, base64.StdEncoding.EncodeToString(image.Data))
		part := imagePart{Type: "image_url"}
		part.ImageUrl.Url = dataUri
		parts = append(parts, part)
	}
	encoded, _ := json.Marshal(parts)
	return encoded
}

func buildTools(tools []llm.Tool) []wireTool {
	if len(tools) == 0 {
		return nil
	}
	wire := make([]wireTool, len(tools))
	for i, tool := range tools {
		wire[i] = wireTool{
			Type: "function",
			Function: wireToolFunc{
				Name:        tool.Name(),
				Description: tool.Description(),
				Parameters:  tool.Schema(),
			},
		}
	}
	return wire
}

func buildResponseFormat(schema *llm.ResponseSchema) *responseFormat {
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

// newChatRequest maps a prompt to a request; Messages are set per iteration.
func newChatRequest(model string, prompt llm.Prompt) chatRequest {
	return chatRequest{
		Model:          model,
		Tools:          buildTools(prompt.Tools),
		ResponseFormat: buildResponseFormat(prompt.ResponseSchema),
	}
}

// finishReasonError maps a non-terminal finish reason to a neutral error. It
// returns nil for the normal reasons ("stop", "tool_calls").
func finishReasonError(reason string) error {
	switch reason {
	case "length":
		return llm.ErrTruncated
	case "content_filter":
		return llm.ErrContentFiltered
	default:
		return nil
	}
}

func parseProviderError(status int, raw []byte) error {
	var response errorResponse
	_ = json.Unmarshal(raw, &response)
	message := response.Error.Message
	if message == "" {
		message = string(raw)
	}
	return &llm.ProviderError{
		Provider:   providerName,
		Type:       response.Error.Type,
		Message:    message,
		StatusCode: status,
	}
}

// mustRawString encodes s as a JSON string; marshaling a string never fails.
func mustRawString(s string) json.RawMessage {
	encoded, _ := json.Marshal(s)
	return encoded
}

// contentString extracts a plain-string content (the final answer). An array
// content is returned verbatim so the caller can still inspect it.
func contentString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var text string
	err := json.Unmarshal(raw, &text)
	if err == nil {
		return text
	}
	return string(raw)
}
