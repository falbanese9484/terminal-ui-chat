package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/falbanese9484/terminal-chat/chat"
	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/providers/models"
	"github.com/falbanese9484/terminal-chat/types"
	"github.com/falbanese9484/terminal-chat/ui"
)

const gap = "\n\n"

/*
This is the main Event Loop for the TUI. It handles initializing a model provider and routing messages to
and from the user.

NOTE: I need to start thinking about UI enhancements, Error Events displayed to the user and
figure out how I'm going to give the user options around the provider and model.
*/
func main() {
	args := os.Args
	var model string
	if len(args) > 1 {
		model = args[1]
	} else {
		model = "x-ai/grok-4-fast:free"
	}
	p := tea.NewProgram(initialModel(model), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg          error
	chatResponseMsg *types.ChatResponse
)

// containing "chat stream closed".
func waitForChatResponse(sub chan *types.ChatResponse) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-sub
		if !ok {
			return errMsg(fmt.Errorf("chat stream closed"))
		}
		return chatResponseMsg(msg)
	}
}

type model struct {
	viewport          viewport.Model
	messages          []string
	textarea          textarea.Model
	debugWindow       viewport.Model
	debugInfo         []string
	showDebug         bool
	senderStyle       lipgloss.Style
	userStyle         lipgloss.Style
	aiStyle           lipgloss.Style
	err               error
	ChatBus           *chat.ChatBus
	ByteReader        chan *types.ChatResponse
	currentAIResponse string
	modelProvider     *types.ProviderService
	modelName         string
	logger            *logger.Logger
	renderer          *glamour.TermRenderer
}

// initialModel creates and returns a fully initialized model configured with a textarea and viewport, styled user and AI label styles, a provider-backed ChatBus with a buffered response channel, a glamour renderer, and a safe logger, and it starts the chat bus goroutine.
func initialModel(m string) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	screenWidth, _, _ := term.GetSize(0)
	mainWidth := screenWidth * 2 / 3
	debugWidth := screenWidth/3 - 4

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	userStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")). // Bright green
		Bold(true)
	aiStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")). // Bright blue
		Bold(true)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	connectedToStyle := lipgloss.NewStyle().Italic(true).
		Foreground(lipgloss.Color("241")) // Gray color for "Connected
	aiConnectedToStyle := connectedToStyle.Foreground(lipgloss.Color("75")) // Lighter blue for AI model name
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("75")). // Bright blue
		Bold(true)
	debugWindowStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Margin(1, 2)
	debugWindowStyle.Width(50)
	debugWindowStyle.Height(10)
	debugView := viewport.New(debugWidth, 10)
	debugView.Style = debugWindowStyle
	debugView.SetContent("Debug Info:\n")
	// Initialize the viewport to be 30 characters wide and 5 characters tall.
	vp := viewport.New(mainWidth, 5)
	ta.SetWidth(mainWidth)
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(2)
	vp.SetContent(logoStyle.Render(ui.LOGO) + titleStyle.Render(ui.PHRASE) + "\n" + connectedToStyle.Render("Connected to: ") + aiConnectedToStyle.Render(m) + "\n\n")
	// Disable newlines in the textarea to handle input on Enter keypress

	ta.KeyMap.InsertNewline.SetEnabled(false)
	logger, err := logger.NewSafeLogger(true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// ollama := models.NewOllamaProvider(logger, m)
	// modelProvider := types.NewProviderService(ollama)
	// TODO: Need to make this part dynamic depending on env or select
	openRouter, err := models.NewOpenRouter(logger, m)
	if err != nil {
		logger.Fatal("failed to initialize openRouter", "error", err)
	}
	modelProvider := types.NewProviderService(openRouter)
	bus := chat.NewChatBus(logger, modelProvider)
	bReader := make(chan *types.ChatResponse, 100)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	go bus.Start(bReader)
	return model{
		textarea:      ta,
		messages:      []string{},
		viewport:      vp,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:           nil,
		ChatBus:       bus,
		userStyle:     userStyle,
		aiStyle:       aiStyle,
		ByteReader:    bReader,
		logger:        logger,
		renderer:      renderer,
		modelProvider: modelProvider,
		modelName:     m,
		debugWindow:   debugView,
		showDebug:     false, // NOTE: Frontend Debugger work in progress via os.Getenv("DEBUG") == "1"
		debugInfo:     []string{},
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

// formatMessage creates a message string prefixed with a timestamped, styled sender label.
// The returned string has the styled "[HH:MM] sender:" prefix followed by a space and the provided content.
func formatMessage(sender, content string, style lipgloss.Style) string {
	timestamp := time.Now().Format("15:04")
	prefix := style.Render(fmt.Sprintf("[%s] %s:", timestamp, sender))
	return prefix + " " + content
}

// and scrolls the viewport to the bottom.
func setAIResponse(m *model, msg *types.ChatResponse) {
	m.currentAIResponse += msg.Response
	renderedText, _ := m.renderer.Render(m.currentAIResponse)
	allMessages := append(m.messages, formatMessage(m.modelName, renderedText, m.aiStyle))
	m.viewport.SetContent(
		lipgloss.NewStyle().Width(
			m.viewport.Width).Render(
			strings.Join(allMessages, "\n")))
	m.viewport.GotoBottom()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Calculate widths
		debugWidth := msg.Width / 3
		mainWidth := msg.Width - debugWidth - 4 // 4 for padding/borders

		// Update main content dimensions
		m.viewport.Width = mainWidth
		m.textarea.SetWidth(mainWidth)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		// Update debug window dimensions
		m.debugWindow.Width = debugWidth
		m.debugWindow.Height = msg.Height

		// Update content
		if len(m.messages) > 0 {
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(
					m.viewport.Width).Render(
					strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case chatResponseMsg:
		if msg == nil {
			// Channel closed without a final message.
			m.logger.Debug("UI:channel closed without a final message")
			return m, nil
		}
		if msg.Response != "" {
			setAIResponse(&m, msg)
		}
		if !msg.Done {
			return m, waitForChatResponse(m.ByteReader)
		} else {
			renderedtext, _ := m.renderer.Render(m.currentAIResponse)
			m.messages = append(m.messages, formatMessage(m.modelName, renderedtext, m.aiStyle))
			m.currentAIResponse = ""
			m.viewport.SetContent(lipgloss.NewStyle().Width(
				m.viewport.Width).Render(
				strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			prompt := m.textarea.Value()
			m.messages = append(m.messages, m.userStyle.Render("You: ")+prompt)
			m.viewport.SetContent(
				lipgloss.NewStyle().Width(
					m.viewport.Width).Render(
					strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			request := m.modelProvider.GenerateRequest(prompt)
			go m.ChatBus.RunChat(request)
			return m, waitForChatResponse(m.ByteReader)
		}

	case errMsg:
		m.err = msg
		m.logger.Debug("UI:frontend error", "error", m.err)
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	// Main content: chat viewport and textarea
	mainContent := fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)

	// When debug is enabled, show it on the right side
	if m.showDebug {
		// Set debug window height to match main content
		m.debugWindow.Height = lipgloss.Height(mainContent)

		// Use lipgloss to join the main content and debug window horizontally
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			mainContent,
			m.debugWindow.View(),
		)
	}

	// When debug is disabled, just show the main content
	return mainContent
}
