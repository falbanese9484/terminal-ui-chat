package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/falbanese9484/terminal-chat/types"
	"github.com/falbanese9484/terminal-chat/ui/styles"
)

func formatMessage(sender, content string, style lipgloss.Style) string {
	timestamp := time.Now().Format("15:04")
	prefix := style.Render(fmt.Sprintf("[%s] %s:", timestamp, sender))
	return prefix + " " + content
}

func setAIResponse(m *ChatModel, msg *types.ChatResponse) {
	m.ChatService.CurrentAIResponse += msg.Response
	renderedText, _ := m.Renderer.Render(m.ChatService.CurrentAIResponse)
	allMessages := append(m.ChatView.Messages, formatMessage(m.ChatService.ModelName, renderedText, styles.AiStyle))
	m.ChatView.Viewport.SetContent(
		lipgloss.NewStyle().Width(
			m.ChatView.Viewport.Width).Render(
			strings.Join(allMessages, "\n")))
	m.ChatView.Viewport.GotoBottom()
}

func waitForChatResponse(sub chan *types.ChatResponse) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-sub
		if !ok {
			return errMsg(fmt.Errorf("chat stream closed"))
		}
		return chatResponsemsg(msg)
	}
}
