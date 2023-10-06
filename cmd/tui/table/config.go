package table

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"winsleepd/cmd/tui/service"
)

func Config() Model {
	columns := []table.Column{
		{Title: "Setting", Width: 25},
		{Title: "Value", Width: 25},
	}

	rows := []table.Row{
		{"timeout", "30m"},
		{"action", "screenoff"},
	}

	service.Get().Configuration = service.GetConfig()
	// convert struct to rows

	//type Configuration struct {
	//	Timeout string `json:"timeout"`
	//	Action  string `json:"action"`
	//}

	rows[0][1] = service.Get().Configuration.Timeout
	rows[1][1] = service.Get().Configuration.Action

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
		Table:       newTable,
		Active:      true,
		TableKeyMap: table.DefaultKeyMap(),
		Additional:  DefaultKeyMap,
		Spinner:     s,
	}
}
