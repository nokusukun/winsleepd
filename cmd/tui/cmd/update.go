package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m DaemonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.TableKeyMap.GotoBottom), key.Matches(msg, m.TableKeyMap.GotoTop):
			m.Table, _ = m.Table.Update(msg)
			return m, nil
		}
	default:
		return m, cmd
	}
	return m, nil
}
