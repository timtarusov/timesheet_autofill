package cmd

import (
	"github.com/spf13/cobra"
	"github.com/timtarusov/timesheet_autofill/pkg"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var YearSend int
var MonthSend int

var send = &cobra.Command{
	Use:   "send",
	Short: "Send an email for reporting",
	Run: func(cmd *cobra.Command, args []string) {

		db, err := gorm.Open(sqlite.Open("history.db"), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}
		pkg.SendEmail(YearSend, MonthSend, db)
	},
}

func init() {
	send.Flags().IntVarP(&YearSend, "year", "y", 2023, "provide year of invoice")
	send.Flags().IntVarP(&MonthSend, "month", "m", 1, "provide month of invoice")
	send.MarkFlagRequired("year")
	send.MarkFlagRequired("month")
	rootCmd.AddCommand(send)
}
