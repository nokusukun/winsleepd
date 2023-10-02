package tui

import "strings"

func (m DaemonModel) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(m.Table.View())
	return stringBuilder.String()
}
