package models

type OpenRouterRequest struct {
	// Request Structure unique to OpenRouter
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}
