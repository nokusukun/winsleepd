package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"winsleepd/cmd/tui/service"
	t "winsleepd/cmd/tui/table"
)

type DaemonModel struct {
	Service *service.Service
	table   t.Model
	keymap  KeyMap
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
		Service: service.Get(),
		table:   t.New(),
		keymap:  DefaultKeyMap,
	}
}

func (m DaemonModel) checkInstall() DaemonModel {
	m.table, _ = m.table.Installed()
	return m
}
