package tui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m DaemonModel) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(m.Table.View())

	return baseStyle.Render(stringBuilder.String())
}
