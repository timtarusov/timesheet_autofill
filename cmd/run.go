package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timtarusov/timesheet_autofill/models"
	"github.com/timtarusov/timesheet_autofill/pkg"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Exclude []int
var Month int
var Year int
var Rate int

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
    db_path := viper.GetString("DB.Path")
		db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}
		db.AutoMigrate(&models.Timesheet{})
		db.AutoMigrate(&models.Invoice{})

		path, err := pkg.Create_path(Year, Month)
		if err != nil {
			log.Fatal(err)
		}
		total := pkg.Write_timesheet(path, Year, Month, Exclude, db)
		pkg.Write_invoice(total, Year, Month, path, Rate, db)
		confirm(path, Year, Month, db)
	},
}

func init() {
	rate := viper.GetInt("Template.Rate")
	run.Flags().IntVarP(&Year, "year", "y", 0, "provide year")
	run.Flags().IntVarP(&Month, "month", "m", 0, "provide month")
	run.Flags().IntSliceVarP(&Exclude, "exclude", "e", []int{}, "provide dates to be excluded")
	run.Flags().IntVarP(&Rate, "rate", "r", rate, "provide rate")
	run.MarkFlagRequired("year")
	run.MarkFlagRequired("month")
	rootCmd.AddCommand(run)
}
