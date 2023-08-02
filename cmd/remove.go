package cmd

import (
	"github.com/spf13/cobra"
	"github.com/timtarusov/timesheet_autofill/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
var YearsR []int
var MonthsR []int

var remove = &cobra.Command{
	Use: "remove",
	Short: "remove db entries",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		db.Where("Year IN ?", YearsR).Where("Month IN ?", MonthsR).Delete(&models.Timesheet{})
		db.Where("Year IN ?", YearsR).Where("Month IN ?", MonthsR).Delete(&models.Invoice{})
	},
}

func init() {
	remove.Flags().IntSliceVarP(&YearsR, "years", "y", []int{}, "years to delete")
	remove.Flags().IntSliceVarP(&MonthsR, "months", "m", []int{}, "months to delete")
	remove.MarkFlagRequired("years")
	remove.MarkFlagRequired("months")
	rootCmd.AddCommand(remove)
}