package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type ChatView struct {
	Viewport viewport.Model
	Messages []string
	renderer *glamour.TermRenderer
}

func NewChatView(width, height int, renderer *glamour.TermRenderer) *ChatView {
	vp := viewport.New(width, height)
	return &ChatView{
		Viewport: vp,
		Messages: []string{},
		renderer: renderer,
	}
}

func (c *ChatView) Set() {
	c.Viewport.SetContent(
		lipgloss.NewStyle().Width(
			c.Viewport.Width).Render(
			strings.Join(c.Messages, "\n")))
	c.Viewport.GotoBottom()
}
