package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/timtarusov/timesheet_autofill/models"
	"github.com/timtarusov/timesheet_autofill/pkg"
)

var (
	Exclude       []int
	Month         int
	Year          int
	Rate          int
	RateInterview int
	NumInterview  int
)

func confirm(path string, year int, month int, db *gorm.DB) {
	prompt := promptui.Prompt{
		Label:     "Proceed?",
		IsConfirm: true,
		Default:   "Y",
	}

	result, err := prompt.Run()

	if (strings.ToLower(result) != "y") && (result != "") {
		fmt.Printf("Removing files at path:  %s\n", path)

		os.RemoveAll(path)
		db.Where("Month=?", month).Where("Year=?", year).Delete(&models.Timesheet{})
		db.Where("Month=?", month).Where("Year=?", year).Delete(&models.Invoice{})
		return
	}
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
}

var run = &cobra.Command{
	Use:   "run",
	Short: "Create timesheet and invoice for the given month",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString("DB.Path")
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		db.AutoMigrate(&models.Timesheet{})
		db.AutoMigrate(&models.Invoice{})

		path, err := pkg.CreatePath(Year, Month)
		if err != nil {
			log.Fatal(err)
		}
		total := pkg.WriteTimesheet(path, Year, Month, Exclude, db)
		pkg.WriteInvoice(total, Year, Month, path, Rate, RateInterview, NumInterview, db)
		confirm(path, Year, Month, db)
	},
}

func init() {
	rate := viper.GetInt("Template.Rate")
	rateInterview := viper.GetInt("Template.RateInterview") // Renamed the local variable to rateInterview
	run.Flags().IntVarP(&Year, "year", "y", 0, "provide year")
	run.Flags().IntVarP(&Month, "month", "m", 0, "provide month")
	run.Flags().IntSliceVarP(&Exclude, "exclude", "e", []int{}, "provide dates to be excluded")
	run.Flags().IntVarP(&Rate, "rate", "r", rate, "provide rate")
	run.Flags().
		IntVarP(&RateInterview, "rateinterview", "v", rateInterview, "provide rate for the interview") // Updated the variable name here
	run.Flags().
		IntVarP(&NumInterview, "numinterview", "i", 0, "provide number of interviews")
	run.MarkFlagRequired("year")
	run.MarkFlagRequired("month")
	rootCmd.AddCommand(run)
}
