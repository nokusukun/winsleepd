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
		refKeyMap := reflect.ValueOf(m.TableKeyMap)
		for i := 0; i < refKeyMap.NumField(); i++ {
			keyMap := refKeyMap.Field(i).Interface().(key.Binding)
			if key.Matches(msg, keyMap) {
				m.Table, cmd = m.Table.Update(msg)
				return m, cmd
			}
		}
		//switch {
		//case key.Matches(msg, m.TableKeyMap.GotoBottom), key.Matches(msg, m.TableKeyMap.GotoTop):
		//	m.Table, _ = m.Table.Update(msg)
		//	return m, nil
		//}
	default:
		return m, cmd
	}
	return m, nil
}
