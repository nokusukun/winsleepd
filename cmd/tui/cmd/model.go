package tui

import (
	"github.com/charmbracelet/bubbles/key"
	t "winsleepd/cmd/tui/table"
)

type DaemonModel struct {
	Table  t.Model
	KeyMap KeyMap
}

type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

var DefaultKeyMap = KeyMap{
	Up:    key.NewBinding(key.WithKeys("pgup", "left", "h")),
	Down:  key.NewBinding(key.WithKeys("pgdown", "right", "l")),
	Enter: key.NewBinding(key.WithKeys("enter")),
	Quit:  key.NewBinding(key.WithKeys("q", "esc", "ctrl+c")),
}

type Emoji struct {
	True, False, Empty string
}

func New() DaemonModel {
	return DaemonModel{
		Table:  t.New(),
		KeyMap: DefaultKeyMap,
	}
}

func (m DaemonModel) running() DaemonModel {
	m.Table = m.Table.Running()
	return m
}
