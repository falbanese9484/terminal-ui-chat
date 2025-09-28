package chat

import (
	"context"

	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
)

type ChatBus struct {
	// This Bus is going to be used to feed messages to the TUI event loop
	Done    chan bool
	Content chan *types.ChatResponse
	// NOTE: ChatResponse will need to be able to take in different response specs.
	Error         chan error
	modelProvider *types.ProviderService
	logger        *logger.Logger
}

/*
I need to start thinking about how I'm going to control different Model specs. I know that that
opencode uses a model registry to give users selection.

My system will most likely be largely based on OpenRouter, with an option to use local models.

TODO: Create Provider Interface that will:
1. Generate a Chat Request to send to the Provider
2. Deserialize the Response into a Chat Response
*/

func NewChatBus(logger *logger.Logger, mp *types.ProviderService) *ChatBus {
	return &ChatBus{
		Done:          make(chan bool),
		Content:       make(chan *types.ChatResponse),
		Error:         make(chan error),
		modelProvider: mp,
		logger:        logger,
	}
}

func (cb *ChatBus) Start(byteReader chan *types.ChatResponse) {
	for {
		select {
		case response := <-cb.Content:
			byteReader <- response
		case err := <-cb.Error:
			cb.logger.Error("failed to read incoming chat response", "error", err)
			return
		case <-cb.Done:
			byteReader <- &types.ChatResponse{Done: true}
		}
	}
}

func (cb *ChatBus) RunChat(request *types.ChatRequest) {
	ctx := context.Background()
	cb.modelProvider.Chat(ctx, request, cb.Content, cb.Error, cb.Done)
}
