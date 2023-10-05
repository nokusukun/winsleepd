package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
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
	Spinner       spinner.Model
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
	False: "❌",
	Empty: "⬜",
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
		{Title: "✨", Width: 6},
		{Title: "Action", Width: 50},
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
		table.WithHeight(len(rows)+2),
	)

	newTable.SetStyles(focused)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		Table:         newTable,
		EnterFunction: functions,
		Active:        true,
		TableKeyMap:   table.DefaultKeyMap(),
		Additional:    DefaultKeyMap,
		Spinner:       s,
	}
}

func (m Model) Install() (Model, tea.Cmd) {
	if !service.Get().IsInstalled() {
		// TODO: Install as user
		service.Get().Install(false)
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

var installed = []table.Row{
	{emoji.True, "Install"},
	{emoji.False, "Start"},
	{emoji.False, "Stop"},
	{emoji.False, "Pause"},
	{emoji.False, "Continue"},
	{emoji.Empty, "Sleep"},
	{emoji.Empty, "Config"},
	{emoji.False, "Debug mode"},
	{emoji.False, "Uninstall"},
}

func (m Model) Installed() (Model, tea.Cmd) {
	if !service.Get().IsInstalled() {
		return New(), nil
	}
	if len(m.Table.Rows()) > 5 {
		return m.Query()
	}
	newRows := installed

	for i, row := range newRows {
		newRows[i][1] = getCellStyle(row[0]).Render(row[1])
	}

	newFuncs := make([]tea.Msg, len(newRows))

	for index, function := range []tea.Msg{
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
		Toggle{},
	} {
		newFuncs[index] = function
	}

	m.Table.SetRows(newRows)
	m.Table.SetHeight(len(newRows) + 1)
	m.EnterFunction = append(m.EnterFunction, newFuncs...)
	return m.Query()
}

func (m Model) Query() (Model, tea.Cmd) {
	if !service.Get().IsInstalled() {
		return New(), nil
	}
	newRows := []table.Row{
		{emoji.True, "Install"},
		{emoji.False, "Start"},
		{emoji.False, "Stop"},
		{emoji.False, "Pause"},
		{emoji.False, "Continue"},
		{emoji.Empty, "Sleep"},
		{emoji.Empty, "Config"},
		{emoji.False, "Debug mode"},
		{emoji.False, "Uninstall"},
	}
	switch service.Get().QueryState() {
	case svc.Running:
		newRows[InstallOpt][0] = emoji.True
		newRows[StartOpt][0] = emoji.True
		newRows[StopOpt][0] = emoji.Empty
		newRows[PauseOpt][0] = emoji.Empty
		newRows[ContinueOpt][0] = emoji.True
	case svc.Stopped:
		newRows[InstallOpt][0] = emoji.True
		newRows[StartOpt][0] = emoji.Empty
		newRows[StopOpt][0] = emoji.True
		newRows[PauseOpt][0] = emoji.False
		newRows[ContinueOpt][0] = emoji.False
	case svc.Paused:
		newRows[InstallOpt][0] = emoji.True
		newRows[StartOpt][0] = emoji.False
		newRows[StopOpt][0] = emoji.Empty
		newRows[PauseOpt][0] = emoji.True
		newRows[ContinueOpt][0] = emoji.Empty
	}

	for i, row := range newRows {
		newRows[i][1] = getCellStyle(row[0]).Render(row[1])
	}

	m.Table.SetRows(newRows)
	return m, nil
}

func getCellStyle(emojiState string) lipgloss.Style {
	switch emojiState {
	case emoji.True, emoji.False:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#6b6b6b"))
	case emoji.Empty:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	}
	// Default style for other cases
	return lipgloss.NewStyle()
}
