// Copyright (c) Liam Stanley <me@liamstanley.io>. All rights reserved. Use
// of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(2)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)
)

type DialogModel struct {
	Id     string
	Height int
	width  int

	Active   string
	Question string
}

func (m DialogModel) Init() tea.Cmd {
	return nil
}

func (m DialogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}

		if zone.Get(m.Id + "confirm").InBounds(msg) {
			m.Active = "confirm"
		} else if zone.Get(m.Id + "cancel").InBounds(msg) {
			m.Active = "cancel"
		}

		return m, nil
	}
	return m, nil
}

func (m DialogModel) View() string {
	var okButton, cancelButton string

	if m.Active == "confirm" {
		okButton = activeButtonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("Maybe")
	} else {
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("Maybe")
	}

	question := lipgloss.NewStyle().Width(27).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		zone.Mark(m.Id+"confirm", okButton),
		zone.Mark(m.Id+"cancel", cancelButton),
	)
	return dialogBoxStyle.Render(lipgloss.JoinVertical(lipgloss.Center, question, buttons))
}
