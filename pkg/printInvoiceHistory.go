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
	totalH := 0.0
	totalUsd := 0.0
	avgRate := 0.0
	const rateFormat = "$ %.2f"

	for _, m := range months {
		row := table.Row{
			m.Year,
			m.Month,
			m.Hours,
			fmt.Sprintf(rateFormat, m.Rate),
			fmt.Sprintf(rateFormat, m.Amount),
		}
		totalH += m.Hours
		totalUsd += m.Amount
		avgRate += m.Rate
		t.AppendRow(row)
	}
	avgRate = avgRate / float64(len(months))
	footer := table.Row{"", "TOTAL/AVG", totalH, fmt.Sprintf(rateFormat, avgRate), fmt.Sprintf(rateFormat, totalUsd)}
	t.AppendFooter(footer)
	t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	t.Render()

}
