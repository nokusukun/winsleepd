package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"winsleepd/cmd/tui/table"
)

func Run() error {
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return err
	}
	return nil
}

func (m DaemonModel) Init() tea.Cmd {
	check := func() tea.Msg {
		return table.Query{}
	}
	return check
}
