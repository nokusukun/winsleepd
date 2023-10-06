package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"reflect"
	"time"
	"winsleepd/cmd/tui/table"
)

func (m main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isInitialized() {
		if _, ok := msg.(tea.WindowSizeMsg); !ok {
			return m, nil
		}
	}
	//var cmd tea.Cmd
	switch msg := msg.(type) {
	case table.Query:
		m.propagate(msg)
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return table.Query{}
		})
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 2
		msg.Width -= 4
		return m.propagate(msg), func() tea.Msg { return table.Query{} } // TODO: Investigate Init() not sending this
	}
	return m.propagate(msg), nil
}

func (m *main) propagate(msg tea.Msg) tea.Model {
	// Propagate to all children.
	m.dialog, _ = m.dialog.Update(msg)
	m.daemon, _ = m.daemon.Update(msg)
	return m
}

func (m DaemonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		msg.Height -= 2
		msg.Width -= 4
		return m, nil
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}
		m.table, _ = m.table.Update(msg)
		return m, nil
	case table.Query:
		m.table, cmd = m.table.Update(msg)
		return m, cmd
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
	}
	return m, nil
}
