package table

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	check := func() tea.Msg {
		return Query{}
	}
	return check
}
