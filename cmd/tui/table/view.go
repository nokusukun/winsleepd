package table

import "strings"

func (m Model) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(m.Table.View())
	stringBuilder.WriteString(m.Spinner.View())
	return stringBuilder.String()
}
