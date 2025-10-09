package components

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/falbanese9484/terminal-chat/logger"
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
		key.WithKeys("esc", "q"),
		key.WithHelp("esc", "cancel"),
	),
}

type ModelItem struct {
	Name string
}

type ModelSelectorCancelMsg struct{}

func (m ModelItem) Title() string       { return m.Name }
func (m ModelItem) Description() string { return "" }
func (m ModelItem) FilterValue() string { return m.Name }

type ModelSelectedMsg struct {
	Name string
}

type ModelSelector struct {
	List         list.Model
	Models       []ModelItem
	Selected     string
	ShowSelector bool
	Width        int
	Height       int
	renderer     *glamour.TermRenderer
	logger       *logger.Logger
}

func NewModelSelector(width, height int, renderer *glamour.TermRenderer, logger *logger.Logger) *ModelSelector {
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
		logger:       logger,
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
	ms.logger.Info("Set models successfully", "items", items)
}

func (ms *ModelSelector) Update(msg tea.Msg) tea.Cmd {
	if !ms.ShowSelector {
		return nil
	}

	// Process the message with the list first, to allow it to handle navigation
	var cmd tea.Cmd
	ms.List, cmd = ms.List.Update(msg)

	// Then handle special cases like Enter and Escape
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		ms.logger.Debug("ModelSelector key pressed", "key", keyMsg.String())

		switch {
		case key.Matches(keyMsg, defaultKeyMap.Select):
			if i, ok := ms.List.SelectedItem().(ModelItem); ok {
				ms.Selected = i.Name
				ms.ShowSelector = false
				return func() tea.Msg {
					return ModelSelectedMsg{Name: i.Name}
				}
			}
		case key.Matches(keyMsg, defaultKeyMap.Cancel):
			ms.ShowSelector = false
			return func() tea.Msg {
				return ModelSelectorCancelMsg{}
			}
		}
	}

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

	if ms.ShowSelector {
		ms.List.ResetSelected()
	}
}
