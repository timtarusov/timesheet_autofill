package pkg

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"github.com/timtarusov/timesheet_autofill/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const (
	SHEET_NAME              = "Sheet1"
	TS_PERIOD_START_DATE    = "B11"
	TS_PERIOD_END_DATE      = "B12"
	TS_CONSULTANT_SIGN_DATE = "K22"
	TS_TOTAL_ROW            = "AG15"
	TS_TOTAL_ALL            = "AG19"
)

func dayInExclude(d int, ex []int) bool {
	for _, i := range ex {
		if i == d {
			return true
		}
	}
	return false
}

func WriteTimesheet(path string, year int, month int, exclude []int, db *gorm.DB) int {
	tsPath := viper.GetString("Template.TimesheetPath")
	tsFn := viper.GetString("Template.TimesheetFilename")
	fmt.Println("Exclude", exclude)
	tsMap := map[int]string{
		1:  "B",
		2:  "C",
		3:  "D",
		4:  "E",
		5:  "F",
		6:  "G",
		7:  "H",
		8:  "I",
		9:  "J",
		10: "K",
		11: "L",
		12: "M",
		13: "N",
		14: "O",
		15: "P",
		16: "Q",
		17: "R",
		18: "S",
		19: "T",
		20: "U",
		21: "V",
		22: "W",
		23: "X",
		24: "Y",
		25: "Z",
		26: "AA",
		27: "AB",
		28: "AC",
		29: "AD",
		30: "AE",
		31: "AF",
	}
	timesheet, err := excelize.OpenFile(tsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := timesheet.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	tsPs := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	tsPe := tsPs.AddDate(0, 1, -1)
	tsSd := tsPs.AddDate(0, 1, 0)
	timesheet.SetCellValue(SHEET_NAME, TS_PERIOD_START_DATE, tsPs)
	timesheet.SetCellValue(SHEET_NAME, TS_PERIOD_END_DATE, tsPe)
	timesheet.SetCellValue(SHEET_NAME, TS_CONSULTANT_SIGN_DATE, tsSd)

	total := 0
	var eights []int
	for i := 1; i < 32; i++ {
		d := time.Date(year, time.Month(month), i, 0, 0, 0, 0, time.Local)
		if d.Month() != time.Month(month) {
			continue
		}
		if (d.Weekday() != time.Saturday) && (d.Weekday() != time.Sunday) && (!dayInExclude(d.Day(), exclude)) {
			cell := fmt.Sprintf("%s%d", tsMap[i], 15)
			cellT := fmt.Sprintf("%s%d", tsMap[i], 19)
			timesheet.SetCellInt(SHEET_NAME, cell, 8)
			timesheet.SetCellInt(SHEET_NAME, cellT, 8)
			total += 8
			eights = append(eights, 8)
		} else {
			eights = append(eights, 0)
		}
	}
	fmt.Printf("In month %d you worked %d days or %d hours\n", month, total/8, total)
	timesheet.SetCellInt(SHEET_NAME, TS_TOTAL_ROW, total)
	timesheet.SetCellInt(SHEET_NAME, TS_TOTAL_ALL, total)

	PrintTSTable(eights, total)
	fmt.Printf("Saving timesheet to %s\n", path+"/"+tsFn)
	if err := timesheet.SaveAs(path + "/" + tsFn); err != nil {
		log.Fatal(err)
	}

	db.Where("Month=?", month).Delete(&models.Timesheet{})
	for i, v := range eights {
		ts := models.Timesheet{
			Day:   i + 1,
			Month: month,
			Year:  year,
			Value: v,
		}
		db.Create(&ts)
	}

	return total
}
