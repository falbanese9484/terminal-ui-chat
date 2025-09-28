package types

type ChatRequest struct {
	// Initial structure for the chat request.
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context,omitempty"`
}

type ChatResponse struct {
	// What we get back from the LLM Api
	Response string `json:"response"`
	Context  []int  `json:"context,omitempty"`
	Done     bool   `json:"done"`
}
