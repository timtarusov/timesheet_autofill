package pkg

import (
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintTSTable(days []int, total int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{}
	eights := table.Row{}
	for i, d := range days {
		header = append(header, strconv.Itoa(i+1))
		if d == 8 {
			eights = append(eights, d)
		} else {
			eights = append(eights, "")
		}
	}
	header = append(header, "TOTAL")
	eights = append(eights, total)
	t.AppendHeader(header)
	t.AppendRow(eights)
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.Render()
}