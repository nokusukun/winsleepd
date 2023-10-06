package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m main) View() string {
	if !m.isInitialized() {
		return ""
	}

	s := baseStyle.MaxHeight(m.height).MaxWidth(m.width).Padding(1, 2, 1, 2)
	return zone.Scan(s.Render(lipgloss.JoinVertical(lipgloss.Center,
		m.dialog.View(), "",
		lipgloss.PlaceHorizontal(
			m.width, lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.daemon.View(),
			),
			lipgloss.WithWhitespaceChars(" "),
		),
	)))
}

func (m DaemonModel) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(m.table.View())
	return baseStyle.Render(stringBuilder.String())
}
