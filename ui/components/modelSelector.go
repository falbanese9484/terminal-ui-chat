package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/falbanese9484/terminal-chat/types"
)

type keyMap struct {
	Select key.Binding
	Cancel key.Binding
}

var defaultKeyMap = keyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select model"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

type ModelItem struct {
	Name string
}

func (m ModelItem) FilterValue() string { return m.Name }

type ModelSelector struct {
	List         list.Model
	Models       []ModelItem
	Selected     string
	ShowSelector bool
	Width        int
	Height       int
	renderer     *glamour.TermRenderer
}

func NewModelSelector(width, height int, renderer *glamour.TermRenderer) *ModelSelector {
	listModel := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	listModel.Title = "Select a Model"
	listModel.SetShowHelp(true)

	return &ModelSelector{
		List:         listModel,
		Models:       []ModelItem{},
		ShowSelector: false,
		Width:        width,
		Height:       height,
		renderer:     renderer,
	}
}

func (ms *ModelSelector) SetModels(models []types.Model) {
	items := []list.Item{}
	for _, model := range models {
		item := ModelItem{
			Name: model.Name,
		}
		items = append(items, item)
	}

	ms.List.SetItems(items)
}

func (ms *ModelSelector) Update(msg tea.Msg) tea.Cmd {
	if !ms.ShowSelector {
		return nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// HandleKeyPress
		switch {
		case key.Matches(msg, defaultKeyMap.Select):
			if i, ok := ms.List.SelectedItem().(ModelItem); ok {
				ms.Selected = i.Name
				ms.ShowSelector = false
				return nil
			}
		case key.Matches(msg, defaultKeyMap.Cancel):
			// Cancel Selection
			ms.ShowSelector = false
			return nil
		}
	}

	ms.List, cmd = ms.List.Update(msg)
	return cmd
}

func (ms *ModelSelector) View() string {
	// If the selector isn't shown, don't render anything
	if !ms.ShowSelector {
		return ""
	}

	// Create a styled container for the list
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")). // Blue border
		Padding(1, 2).
		Width(ms.Width).
		Align(lipgloss.Center)

	// Return the styled list view
	return style.Render(ms.List.View())
}

func (ms *ModelSelector) Toggle() {
	ms.ShowSelector = !ms.ShowSelector

	// When showing the selector, ensure the list has focus
	if ms.ShowSelector {
		// You might want to refresh models here or in the caller
	}
}
