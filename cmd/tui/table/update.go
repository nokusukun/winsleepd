package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"winsleepd/cmd/tui/service"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Toggle:
		index := m.Table.Cursor()
		switch index {
		case InstallOpt:
			return m.Install()
		case StartOpt:
			return m.Start()
		case StopOpt:
			return m.Stop()
		case PauseOpt:
			return m.Pause()
		case ContinueOpt:
			return m.Continue()
		case SleepOpt:
			service.Get().Sleep()
		case ConfigOpt:
			service.Get().Config()
		case DebugOpt:
			// TODO: Implement debug mode
		case UninstallOpt:
			return m.Uninstall()
		}
		return m.Query()
	case CheckInstalled:
		return m.Installed()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Additional.Enter):
			if ok := m.EnterFunction[m.Table.Cursor()]; ok != nil {
				//m.Table.Blur()
				return m.Update(ok)
			}

		}
		m.Table, _ = m.Table.Update(msg)
		return m, nil
	}
	return m, nil
}

type CheckInstalled struct{}
type Install struct{}
type Toggle struct{}

const (
	InstallOpt = iota
	StartOpt
	StopOpt
	PauseOpt
	ContinueOpt
	SleepOpt
	ConfigOpt
	DebugOpt
	UninstallOpt
)
