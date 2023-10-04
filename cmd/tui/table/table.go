package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"winsleepd/cmd/tui/service"
)

type Model struct {
	Table         table.Model
	Service       *service.Service
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
	True:  "✔️",
	False: "",
	Empty: "",
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
	}

	functions := make([]tea.Msg, len(rows))

	for index, function := range []tea.Msg{
		Install{},
	} {
		functions[index] = function
	}

	newTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(8),
	)

	newTable.SetStyles(focused)

	return Model{
		Service:       service.Get(),
		Table:         newTable,
		EnterFunction: functions,
		Active:        true,
		TableKeyMap:   table.DefaultKeyMap(),
		Additional:    DefaultKeyMap,
	}
}

func (m Model) Install() (Model, tea.Cmd) {
	if !service.Get().IsInstalled() {
		service.Get().Install()
	}
	return m.Running()
}

func (m Model) Running() (Model, tea.Cmd) {
	if !service.Get().IsInstalled() {
		return m, nil
	}
	currentRows := m.Table.Rows()
	if len(currentRows) > 5 {
		return m, nil
	}
	newRows := []table.Row{
		{emoji.False, "Start"},
		{emoji.False, "Stop"},
		{emoji.False, "Pause"},
		{emoji.Empty, "Sleep"},
		{emoji.Empty, "Config"},
		{emoji.False, "Debug mode"},
	}

	newFuncs := make([]tea.Msg, len(newRows))

	for index, function := range []tea.Msg{
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
	} {
		newFuncs[index] = function
	}

	rows := append(currentRows, newRows...)
	m.Table.SetRows(rows)
	m.EnterFunction = append(m.EnterFunction, newFuncs...)
	return m.Update(Toggle{})
}
