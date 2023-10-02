package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func Run() error {
	p := tea.NewProgram(New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		return err
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}
