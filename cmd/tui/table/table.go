package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Table         table.Model
	EnterFunction []tea.Msg
	Active        bool
	TableKeyMap   table.KeyMap
	Additional    KeyMap
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

var emoji = Emoji{
	True:  "✅",
	False: "⬜",
	Empty: "⬛",
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
		{emoji.False, "Install"},
		{emoji.False, "Start"},
		{emoji.False, "Stop"},
		{emoji.False, "Pause"},
		{emoji.Empty, "Sleep"},
	}

	functions := make([]tea.Msg, len(rows))

	f := []tea.Msg{
		nil,
	}

	for i := range f {
		functions[i] = f[i]
	}

	newTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(4),
	)

	newTable.SetStyles(focused)

	return Model{
		Table:         newTable,
		EnterFunction: functions,
		Active:        true,
		TableKeyMap:   table.DefaultKeyMap(),
		Additional:    DefaultKeyMap,
	}
}

func (m Model) Running() Model {
	currentRows := m.Table.Rows()
	newRows := []table.Row{
		{emoji.Empty, "Config"},
		{emoji.False, "Debug mode"},
	}
	rows := append(currentRows, newRows...)
	m.Table.SetRows(rows)
	return m
}
