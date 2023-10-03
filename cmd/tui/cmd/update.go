package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
)

func (m DaemonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		refKeyMap := reflect.ValueOf(m.Table.TableKeyMap)
		for i := 0; i < refKeyMap.NumField(); i++ {
			keyMap := refKeyMap.Field(i).Interface().(key.Binding)
			if key.Matches(msg, keyMap) {
				newTable, _ := m.Table.Update(msg)
				m.Table = newTable
				return m, nil
			}
		}
		switch {
		case key.Matches(msg, m.KeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.Enter):
			newTable, _ := m.Table.Update(msg)
			m.Table = newTable
			return m, nil
		}
	default:
		return m, cmd
	}
	return m, nil
}
