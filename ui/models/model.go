package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
	"github.com/falbanese9484/terminal-chat/ui/components"
	"github.com/falbanese9484/terminal-chat/ui/services"
	"github.com/falbanese9484/terminal-chat/ui/styles"
)

type (
	chatResponsemsg *types.ChatResponse
	errMsg          error
	UIMode          int
)

const (
	gap             = "\n\n"
	ChatMode UIMode = iota
	ModelSelectMode
)

type ChatModel struct {
	InputArea     *components.InputArea
	ChatView      *components.ChatView
	DebugView     *components.DebugWindow
	ModelSelector *components.ModelSelector
	ChatService   *services.ChatService
	Logger        *logger.Logger
	Renderer      *glamour.TermRenderer
	Err           error
	Mode          UIMode
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m ChatModel) handleChatResponse(msg chatResponsemsg) (tea.Model, tea.Cmd) {
	if msg == nil {
		m.Logger.Debug("UI:channel closed without a final message")
		return m, nil
	}
	if msg.Response != "" {
		setAIResponse(&m, msg)
	}
	if !msg.Done {
		return m, waitForChatResponse(m.ChatService.ByteReader)
	} else {
		renderedText, _ := m.Renderer.Render(m.ChatService.CurrentAIResponse)
		m.ChatView.Messages = append(m.ChatView.Messages, formatMessage(m.ChatService.ModelName, renderedText, styles.AiStyle))
		m.ChatView.Set()
	}
	return m, nil
}

func (m ChatModel) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	debugWidth := msg.Width / 3
	mainWidth := msg.Width - debugWidth - 4

	m.ChatView.Viewport.Width = mainWidth
	m.InputArea.Textarea.SetWidth(mainWidth)
	m.ChatView.Viewport.Height = msg.Height - m.InputArea.Textarea.Height() - lipgloss.Height(gap)

	m.DebugView.Viewport.Width = debugWidth
	m.DebugView.Viewport.Height = msg.Height

	if len(m.ChatView.Messages) > 0 {
		m.ChatView.Set()
	}
	return m, nil
}

func (m ChatModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEscape:
		if m.Mode == ModelSelectMode {
			m.ModelSelector.Toggle()
			m.Mode = ChatMode
			return m, nil
		}
		fmt.Println(m.InputArea.Textarea.Value())
		return m, tea.Quit
		//	case tea.KeyCtrlM:
		//		m.ModelSelector.Toggle()
		//		if m.ModelSelector.ShowSelector {
		//			m.Mode = ModelSelectMode
		//		} else {
		//			m.Mode = ChatMode
		//		}
		//		return m, nil //TODO: Key binding issue - cant use ctrlM and Enter in the same statement
	case tea.KeyEnter:
		prompt := m.InputArea.Textarea.Value()
		m.ChatView.Messages = append(m.ChatView.Messages, styles.UserStyle.Render("You: ")+prompt)
		m.ChatView.Set()
		m.InputArea.Textarea.Reset()
		request := m.ChatService.ModelProvider.GenerateRequest(prompt)
		go m.ChatService.Bus.RunChat(request)
		return m, waitForChatResponse(m.ChatService.ByteReader)
	}
	return m, nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.InputArea.Textarea, tiCmd = m.InputArea.Textarea.Update(msg)
	m.ChatView.Viewport, vpCmd = m.ChatView.Viewport.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)
	case chatResponsemsg:
		return m.handleChatResponse(msg)
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case errMsg:
		m.Err = msg
		m.Logger.Debug("UI:frontend error", "error", m.Err)
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	mainContent := fmt.Sprintf(
		"%s%s%s",
		m.ChatView.Viewport.View(),
		gap,
		m.InputArea.Textarea.View(),
	)

	if m.DebugView.ShowDebug {
		m.DebugView.Viewport.Height = lipgloss.Height(mainContent)
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			mainContent,
			m.DebugView.Viewport.View(),
		)
	}

	return mainContent
}
