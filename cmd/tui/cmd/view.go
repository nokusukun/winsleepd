package tui

import "strings"

func (m Model) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(m.Table.View())
	return stringBuilder.String()
}
