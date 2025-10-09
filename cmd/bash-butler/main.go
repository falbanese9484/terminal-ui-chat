package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/x/term"
	"github.com/falbanese9484/terminal-chat/chat"
	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/providers/models"
	"github.com/falbanese9484/terminal-chat/types"
	"github.com/falbanese9484/terminal-chat/ui"
	"github.com/falbanese9484/terminal-chat/ui/components"
	uiModels "github.com/falbanese9484/terminal-chat/ui/models"
	"github.com/falbanese9484/terminal-chat/ui/services"
	"github.com/falbanese9484/terminal-chat/ui/styles"
)

func main() {
	args := os.Args
	var modelName string
	if len(args) > 1 {
		modelName = args[1]
	} else {
		modelName = "llama3.2:latest"
	}
	p := tea.NewProgram(initialModel(modelName), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func initialModel(modelName string) tea.Model {
	// Get screen dimensions
	screenWidth, _, _ := term.GetSize(0)
	mainWidth := screenWidth * 2 / 3

	// Initialize logger
	logger, err := logger.NewSafeLogger(true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Initialize glamour renderer
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	// initialize the ModelRefresher
	modelRefresher := types.NewModelRefresher(3600)

	// Initialize provider
	// openRouter, err := models.NewOpenRouter(logger, modelName, modelRefresher)
	ollama := models.NewOllamaProvider(logger, modelName, modelRefresher)
	//if err != nil {
	//	logger.Fatal("failed to initialize openRouter", "error", err)
	//}
	modelProvider := types.NewProviderService(ollama)

	// Initialize chat bus and response channel
	bus := chat.NewChatBus(logger, modelProvider)
	byteReader := make(chan *types.ChatResponse, 100)

	// Create chat service
	chatService := &services.ChatService{
		Bus:               bus,
		ByteReader:        byteReader,
		CurrentAIResponse: "",
		ModelProvider:     modelProvider,
		ModelName:         modelName,
		Logger:            logger,
	}

	// Create UI components
	inputArea := components.NewInputArea(renderer)
	chatView := components.NewChatView(mainWidth, screenWidth/2, renderer)
	modelSelector := components.NewModelSelector(mainWidth, screenWidth/8, renderer, logger)

	// Initialize chatView with logo and connection message
	logoContent := styles.LogoStyle.Render(ui.LOGO) + styles.TitleStyle.Render(ui.PHRASE) + "\n" +
		styles.ConnectedToStyle.Render("Connected to: ") + styles.AiConnectedToStyle.Render(modelName) + "\n\n"
	chatView.Viewport.SetContent(logoContent)

	// Start chat bus
	go bus.Start(byteReader)

	// Create and return the chat model
	return &uiModels.ChatModel{
		InputArea:     inputArea,
		ChatView:      chatView,
		ChatService:   chatService,
		Logger:        logger,
		Renderer:      renderer,
		Err:           nil,
		ModelSelector: modelSelector,
	}
}
