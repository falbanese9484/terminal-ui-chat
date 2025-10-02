package styles

import "github.com/charmbracelet/lipgloss"

var (
	UserStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	AiStyle   = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")). // Bright blue
			Bold(true)
	ConnectedToStyle = lipgloss.NewStyle().Italic(true).
				Foreground(lipgloss.Color("241")) // Gray color for "Connected
	AiConnectedToStyle = ConnectedToStyle.Foreground(lipgloss.Color("75")) // Lighter blue for AI model name
	LogoStyle          = lipgloss.NewStyle().
				Foreground(lipgloss.Color("75")). // Bright blue
				Bold(true)
	DebugWindowStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1, 2).
				Margin(1, 2)
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(2)
)
