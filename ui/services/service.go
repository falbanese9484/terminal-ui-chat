package services

import (
	"github.com/falbanese9484/terminal-chat/chat"
	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
)

type ChatService struct {
	Bus               *chat.ChatBus
	ByteReader        chan *types.ChatResponse
	CurrentAIResponse string
	ModelProvider     *types.ProviderService
	ModelName         string
	Logger            *logger.Logger
}
