package pkg

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/timtarusov/timesheet_autofill/models"
)

func PrintInvoiceHistory(months []*models.Invoice) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	header := table.Row{
		"Year", "Month", "Hours", "Rate", "Amount",
	}
	t.AppendHeader(header)
	total_h := 0.0
	total_usd := 0.0
	avg_rate := 0.0
	for _, m := range months {
		row := table.Row{
			m.Year,
			m.Month,
			m.Hours,
			fmt.Sprintf("$ %.2f", m.Rate),
			fmt.Sprintf("$ %.2f", m.Amount),
		}
		total_h += m.Hours
		total_usd += m.Amount
		avg_rate += m.Rate
		t.AppendRow(row)
	}
	avg_rate = avg_rate/float64(len(months))
	footer := table.Row{"", "TOTAL/AVG", total_h, fmt.Sprintf("$ %.2f", avg_rate), fmt.Sprintf("$ %.2f", total_usd)}
	t.AppendFooter(footer)
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.Render()

}