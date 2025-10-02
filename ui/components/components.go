package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ChatView struct {
	Viewport viewport.Model
	Messages []string
	renderer *glamour.TermRenderer
}

type InputArea struct {
	Textarea textarea.Model
	renderer *glamour.TermRenderer
}

type DebugWindow struct {
	Viewport  viewport.Model
	logs      []string
	ShowDebug bool
	renderer  *glamour.TermRenderer
}

func NewInputArea(renderer *glamour.TermRenderer) *InputArea {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	return &InputArea{
		Textarea: ta,
		renderer: renderer,
	}
}

func NewChatView(width, height int, renderer *glamour.TermRenderer) *ChatView {
	vp := viewport.New(width, height)
	return &ChatView{
		Viewport: vp,
		Messages: []string{},
		renderer: renderer,
	}
}

func NewDebugWindow(width, height int, show bool, renderer *glamour.TermRenderer) *DebugWindow {
	vp := viewport.New(width, height)
	return &DebugWindow{
		Viewport:  vp,
		logs:      []string{},
		ShowDebug: show,
		renderer:  renderer,
	}
}

func (c *ChatView) Set() {
	c.Viewport.SetContent(
		lipgloss.NewStyle().Width(
			c.Viewport.Width).Render(
			strings.Join(c.Messages, "\n")))
	c.Viewport.GotoBottom()
}
