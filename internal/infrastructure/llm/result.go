package llm

type Result struct {
	Content string
	Usage   Usage
}

func NewResult(content string, usage Usage) *Result {
	return &Result{
		Content: content,
		Usage:   usage,
	}
}

// Usage is summed across the tool-call loop; Calls is the number of round trips.
type Usage struct {
	InputTokens  int
	OutputTokens int
	Calls        int
}

// AddCall records one provider round trip's token usage.
func (u *Usage) AddCall(inputTokens, outputTokens int) {
	u.InputTokens += inputTokens
	u.OutputTokens += outputTokens
	u.Calls++
}
