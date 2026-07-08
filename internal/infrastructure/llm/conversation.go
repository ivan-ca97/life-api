package llm

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type Image struct {
	MediaType string // e.g. "image/jpeg"
	Data      []byte
}

type Message struct {
	Role   Role
	Text   string
	Images []Image
}

type Conversation struct {
	Instructions string
	Messages     []Message
}

func SingleTurn(instructions, text string, images []Image) Conversation {
	message := Message{
		Role:   RoleUser,
		Text:   text,
		Images: images,
	}
	return Conversation{
		Instructions: instructions,
		Messages:     []Message{message},
	}
}
