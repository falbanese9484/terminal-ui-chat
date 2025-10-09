package components

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type InputArea struct {
	Textarea textarea.Model
	renderer *glamour.TermRenderer
}

func NewInputArea(renderer *glamour.TermRenderer) *InputArea {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "| "
	ta.SetHeight(3)

	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	return &InputArea{
		Textarea: ta,
		renderer: renderer,
	}
}
