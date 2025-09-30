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

// NewChatBus creates and returns a ChatBus with initialized channels for Done, Content, and Error,
// and assigns the provided logger and model provider.
// The channels are unbuffered.

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
			cb.logger.Debug("streaming", "bytes", len(cb.Content))
			byteReader <- response
		case err := <-cb.Error:
			cb.logger.Error("failed to read incoming chat response", "error", err)
			return
		case <-cb.Done:
			cb.logger.Info("message complete - signalling done")
			byteReader <- &types.ChatResponse{Done: true}
		}
	}
}

func (cb *ChatBus) RunChat(request *types.ChatRequest) {
	ctx := context.Background()
	cb.modelProvider.Chat(&types.BusConnector{Ctx: ctx, Request: request, ResponseChan: cb.Content, ErrorChan: cb.Error, DoneChannel: cb.Done})
}
