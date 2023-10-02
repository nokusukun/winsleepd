package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Table       table.Model
	Active      bool
	TableKeyMap table.KeyMap
	KeyMap      KeyMap
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

var focused = table.DefaultStyles()
var unfocused = table.DefaultStyles()

func New() Model {
	focused = table.DefaultStyles()
	focused.Header = focused.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	focused.Selected = focused.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	unfocused = table.DefaultStyles()
	unfocused.Header = unfocused.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	unfocused.Selected = unfocused.Selected.
		Foreground(lipgloss.Color("#bbbbbb"))

	columns := []table.Column{
		{Title: "#", Width: 6},
		{Title: "Stats", Width: 50},
	}

	rows := []table.Row{
		{"0", "Images with captions"},
		{"0", "Images with captions that match directories"},
		{"0", "Missing captions"},
		{"0", "Pending text files"},
	}

	return Model{
		Table: table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(4),
		),
		Active:      true,
		TableKeyMap: table.DefaultKeyMap(),
		KeyMap:      DefaultKeyMap,
	}
}
