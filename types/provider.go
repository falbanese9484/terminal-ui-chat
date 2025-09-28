package types

import (
	"context"
)

type Provider interface {
	// Provider Takes the universal Chat Request, along with reponseChannel and errorChannel
	// from the Chat Bus. This way each Provider can handle the serializing, deserializing
	// and streaming that may be provider specific.
	Chat(ctx context.Context, request *ChatRequest,
		responseChannel chan *ChatResponse,
		errorChannel chan error, doneChannel chan bool)
	GenerateRequest(prompt string) *ChatRequest
}

type ProviderService struct {
	modelProvider Provider
}

// NewProviderService creates a ProviderService that uses the given Provider as its modelProvider.
func NewProviderService(mp Provider) *ProviderService {
	return &ProviderService{
		modelProvider: mp,
	}
}

func (ps *ProviderService) Chat(ctx context.Context, request *ChatRequest,
	responseChannel chan *ChatResponse,
	errorChannel chan error, doneChannel chan bool,
) {
	ps.modelProvider.Chat(ctx, request, responseChannel, errorChannel, doneChannel)
}

func (ps *ProviderService) GenerateRequest(prompt string) *ChatRequest {
	return ps.modelProvider.GenerateRequest(prompt)
}
