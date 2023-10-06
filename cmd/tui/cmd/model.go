package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"winsleepd/cmd/tui/dialog"
	"winsleepd/cmd/tui/service"
	t "winsleepd/cmd/tui/table"
)

type main struct {
	height int
	width  int

	daemon tea.Model
	dialog tea.Model
}

func (m main) isInitialized() bool {
	return m.height != 0 && m.width != 0
}

type DaemonModel struct {
	Service *service.Service
	table   t.Model
	keymap  KeyMap
	width   int
	height  int
	Id      string
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

func NewMain() *main {
	return &main{
		height: 0,
		width:  0,
		daemon: *newDaemon(),
		dialog: *newDialog(),
	}
}

func newDialog() *dialog.DialogModel {
	return &dialog.DialogModel{
		Id:       zone.NewPrefix(),
		Height:   8,
		Active:   "confirm",
		Question: "Are you sure you want to eat marmalade?",
	}
}

func newDaemon() *DaemonModel {
	return &DaemonModel{
		Id:      zone.NewPrefix(),
		Service: service.Get(),
		table:   t.New(),
		keymap:  DefaultKeyMap,
		width:   0,
		height:  0,
	}
}

func (m DaemonModel) checkInstall() DaemonModel {
	m.table, _ = m.table.Installed()
	return m
}
