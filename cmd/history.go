package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timtarusov/timesheet_autofill/models"
	"github.com/timtarusov/timesheet_autofill/pkg"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Months []int
	Years  []int
)

var history = &cobra.Command{
	Use:   "history",
	Short: "Show history of invoices and timesheets",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString("DB.Path")
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		err = db.AutoMigrate(&models.Timesheet{})
		if err != nil {
			panic("failed to migrate")
		}
		err = db.AutoMigrate(&models.Invoice{})
		if err != nil {
			panic("failed to migrate")
		}

		var TS []*models.Timesheet
		var IN []*models.Invoice

		const orderAsc = "year asc"

		const orderMonthAsc = "month asc"

		if len(Months) == 0 {
			db.Order(orderAsc).Order(orderMonthAsc).Order("day asc").Find(&TS)
			db.Order(orderAsc).Order(orderMonthAsc).Find(&IN)
		} else {
			var cond []string
			for _, m := range Months {
				cond = append(cond, strconv.Itoa(m))
			}

			db.Order(orderAsc).Order(orderMonthAsc).Order("day asc").Find(&TS, "month IN ?", cond)
			db.Order(orderAsc).Order(orderMonthAsc).Find(&IN, "month IN ?", cond)
		}
		fmt.Println("Timesheets for the chosen months:")
		pkg.PrintTSHistory(TS)
		fmt.Print("\n")
		fmt.Println("Invoices for the chosen months:")
		pkg.PrintInvoiceHistory(IN)
	},
}

func init() {
	history.Flags().IntSliceVarP(&Months, "months", "m", []int{}, "Provide months for search")
	history.Flags().IntSliceVarP(&Years, "years", "y", []int{}, "Provide years for search")
	rootCmd.AddCommand(history)
}
