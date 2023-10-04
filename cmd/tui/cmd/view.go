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
	stringBuilder.WriteString(m.table.View())
	return baseStyle.Render(stringBuilder.String())
}
