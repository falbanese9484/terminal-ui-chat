package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/joho/godotenv/autoload"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/falbanese9484/terminal-chat/chat"
	"github.com/falbanese9484/terminal-chat/logger"
)

const gap = "\n\n"

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg          error
	chatResponseMsg *chat.ChatResponse
)

func waitForChatResponse(sub chan *chat.ChatResponse) tea.Cmd {
	return func() tea.Msg {
		return chatResponseMsg(<-sub)
	}
}

type model struct {
	viewport          viewport.Model
	messages          []string
	textarea          textarea.Model
	senderStyle       lipgloss.Style
	err               error
	ChatBus           *chat.ChatBus
	ByteReader        chan *chat.ChatResponse
	currentAIResponse string
	logger            *logger.Logger
	context           []int
	renderer          *glamour.TermRenderer
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)
	logger, err := logger.NewSafeLogger(true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	bus := chat.NewChatBus(logger)
	bReader := make(chan *chat.ChatResponse, 100)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	go bus.Start(bReader)
	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		ChatBus:     bus,
		ByteReader:  bReader,
		logger:      logger,
		context:     []int{},
		renderer:    renderer,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
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
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case chatResponseMsg:
		if msg.Response != "" {
			m.currentAIResponse += msg.Response
			renderedText, _ := m.renderer.Render(m.currentAIResponse)
			allMessages := append(m.messages, m.senderStyle.Render("AI: ")+renderedText)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(allMessages, "\n")))
			m.viewport.GotoBottom()
		}
		if !msg.Done {
			return m, waitForChatResponse(m.ByteReader)
		} else {
			renderedtext, _ := m.renderer.Render(m.currentAIResponse)
			m.messages = append(m.messages, m.senderStyle.Render("AI: ")+renderedtext)
			m.currentAIResponse = ""
			m.context = msg.Context
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			prompt := m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+prompt)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			request := chat.ChatRequest{
				Model:   "llama3.2",
				Prompt:  prompt,
				Stream:  true,
				Context: m.context,
			}
			go m.ChatBus.RunChat(&request)
			return m, waitForChatResponse(m.ByteReader)
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
