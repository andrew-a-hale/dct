package utils

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var style = lipgloss.NewStyle().Align(lipgloss.Center)

func Render(writer io.Writer, result Result, maxRows int) error {
	var headers []string
	var types []string
	for _, header := range result.Headers {
		headers = append(headers, header.Name)
		types = append(types, header.Type)
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case 3:
				// force border after type row in display
				return style.
					Border(lipgloss.NormalBorder(), true, false, false, false)
			default:
				return style
			}
		}).
		Rows(headers).
		Rows(types)

	rowsToDisplay := min(maxRows, len(result.Rows))

	t.Rows(result.Rows[:rowsToDisplay]...)

	fmt.Println(t)
	return nil
}
