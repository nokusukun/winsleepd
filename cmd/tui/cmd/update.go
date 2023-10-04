package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
	"winsleepd/cmd/tui/service"
)

func (m DaemonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case service.Service:
		m.table.Update(msg)
	case tea.KeyMsg:
		refKeyMap := reflect.ValueOf(m.table.TableKeyMap)
		for i := 0; i < refKeyMap.NumField(); i++ {
			keyMap := refKeyMap.Field(i).Interface().(key.Binding)
			if key.Matches(msg, keyMap) {
				newTable, _ := m.table.Update(msg)
				m.table = newTable
				return m, nil
			}
		}
		switch {
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.Enter):
			newTable, _ := m.table.Update(msg)
			m.table = newTable
			return m, nil
		}
	default:
		return m, cmd
	}
	return m, nil
}
