package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/sys/windows/svc"
	"winsleepd/cmd/tui/service"
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
		Toggle{},
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
	return m.Installed()
}

func (m Model) Start() (Model, tea.Cmd) {
	service.Get().Start()
	return m.Query()
}

func (m Model) Stop() (Model, tea.Cmd) {
	service.Get().Stop()
	return m.Query()
}

func (m Model) Pause() (Model, tea.Cmd) {
	service.Get().Pause()
	return m.Query()
}

func (m Model) Continue() (Model, tea.Cmd) {
	service.Get().Continue()
	return m.Query()
}

func (m Model) Uninstall() (Model, tea.Cmd) {
	service.Get().Uninstall()
	return m.Installed()
}

func (m Model) Installed() (Model, tea.Cmd) {
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
		{emoji.False, "Continue"},
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
	return m.Query()
}

func (m Model) Query() (Model, tea.Cmd) {
	currentRows := m.Table.Rows()
	if len(currentRows) > 5 {
		return m, nil
	}
	switch service.Get().QueryState() {
	case svc.Running:
		currentRows[StartOpt][0] = emoji.True
		currentRows[StopOpt][0] = emoji.False
		currentRows[PauseOpt][0] = emoji.False
		currentRows[ContinueOpt][0] = emoji.True
	case svc.Stopped:
		currentRows[StartOpt][0] = emoji.False
		currentRows[StopOpt][0] = emoji.True
		currentRows[PauseOpt][0] = emoji.False
		currentRows[ContinueOpt][0] = emoji.False
	case svc.Paused:
		currentRows[StartOpt][0] = emoji.False
		currentRows[StopOpt][0] = emoji.False
		currentRows[PauseOpt][0] = emoji.True
		currentRows[ContinueOpt][0] = emoji.False
	}
	m.Table.SetRows(currentRows)
	return m, nil
}
