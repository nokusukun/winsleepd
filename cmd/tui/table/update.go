package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Toggle:
		currentRows := m.Table.Rows()
		index := m.Table.Cursor()
		switch m.Table.SelectedRow()[0] {
		case emoji.True:
			currentRows[index][0] = emoji.False
		case emoji.False:
			currentRows[index][0] = emoji.True
		}
		m.Table.SetRows(currentRows)
		return m, nil
	case CheckIfRunning:
		return m.Running()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Additional.Enter):
			if ok := m.EnterFunction[m.Table.Cursor()]; ok != nil {
				//m.Table.Blur()
				return m.Update(ok)
			}

		}
		m.Table, _ = m.Table.Update(msg)
		return m, nil
	}
	return m, nil
}

type CheckIfRunning struct{}
type Toggle struct{}
