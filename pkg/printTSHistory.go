package pkg

import (
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/timtarusov/timesheet_autofill/models"
)

func PrintTSHistory(months []*models.Timesheet) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{}
	header = append(header, "YEAR")
	header = append(header, "MONTH")
	for i := 1; i < 32; i++ {
		header = append(header, strconv.Itoa(i))
	}
	header = append(header, "TOTAL")
	t.AppendHeader(header)

	var monthsMap = make(map[int][]*models.Timesheet)
	for _, m := range months {
		monthsMap[m.Month] = append(monthsMap[m.Month], m)
	}
	for m, v := range monthsMap {
		row := table.Row{}
		total := 0
		row = append(row, 0)
		row = append(row, m)
		for _, d := range v {
			row[0] = d.Year
			if d.Value == 8 {
				row = append(row, d.Value)
				total += d.Value
			} else {
				row = append(row, "")
			}
		}
		diff := 31 - len(v)
		for i := 0; i < diff; i++ {
			row = append(row, "")
		}
		row = append(row, total)
		t.AppendRow(row)
	}
	t.SortBy([]table.SortBy{
		{Name: "YEAR", Mode: table.AscNumeric},
		{Name: "MONTH", Mode: table.AscNumeric},
	})
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.Render()
}
