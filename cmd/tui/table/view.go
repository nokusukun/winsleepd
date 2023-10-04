package table

import "strings"

func (m Model) View() string {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(m.Spinner.View())
	stringBuilder.WriteString(m.Table.View())
	return stringBuilder.String()
}
