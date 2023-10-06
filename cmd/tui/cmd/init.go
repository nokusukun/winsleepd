package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"winsleepd/cmd/tui/table"
)

func Run() error {
	zone.NewGlobal()
	p := tea.NewProgram(NewMain(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return err
	}
	return nil
}

func (m main) Init() tea.Cmd {
	check := func() tea.Msg {
		return table.Query{}
	}
	return check
}

func (m DaemonModel) Init() tea.Cmd {
	check := func() tea.Msg {
		return table.Query{}
	}
	return check
}
