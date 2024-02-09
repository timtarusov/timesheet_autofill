package pkg

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"github.com/timtarusov/timesheet_autofill/models"
)

const (
	IN_SHEET             = "Simple Invoice"
	IN_DATE              = "E4"
	IN_NUM               = "E5"
	IN_DESC              = "A20"
	IN_UNITS             = "C20"
	IN_RATE              = "D20"
	IN_AMOUNT            = "E20"
	IN_TOTAL             = "E37"
	IN_DESC_INTERVIEWS   = "A21"
	IN_UNITS_INTERVIEWS  = "C21"
	IN_RATE_INTERVIEWS   = "D21"
	IN_AMOUNT_INTERVIEWS = "E21"
	TIME_FORMAT          = "Jan 2, 2006"
)

func WriteInvoice(
	t int,
	year int,
	month int,
	path string,
	rate int,
	rateInt int,
	nInt int,
	db *gorm.DB,
) int {
	templInvoice := viper.GetString("Template.InvoicePath")
	invFn := viper.GetString("Template.InvoiceFilename")
	inv, err := excelize.OpenFile(templInvoice)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := inv.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	d := time.Date(year, time.Month(month), 15, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	num, err := inv.GetCellValue(IN_SHEET, IN_NUM)
	if err != nil {
		log.Fatal(err)
	}
	numInt, err := strconv.Atoi(num)
	if err != nil {
		log.Fatal(err)
	}

	totalUsd := t * rate
	totalInterviews := nInt * rateInt
	totalOverall := totalUsd + totalInterviews

	descDS := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	descDE := descDS.AddDate(0, 1, -1)
	desc := fmt.Sprintf(
		"Consulting Servives - for the period of %s - %s",
		descDS.Format(TIME_FORMAT),
		descDE.Format(TIME_FORMAT),
	)
	descInterview := fmt.Sprintf(
		"Technical interviews - for the period of %s - %s",
		descDS.Format(TIME_FORMAT),
		descDE.Format(TIME_FORMAT),
	)

	inv.SetCellValue(IN_SHEET, IN_DATE, d)
	inv.SetCellValue(IN_SHEET, IN_NUM, numInt+1)
	inv.SetCellInt(IN_SHEET, IN_UNITS, t)
	inv.SetCellValue(IN_SHEET, IN_DESC, desc)
	inv.SetCellValue(IN_SHEET, IN_RATE, rate)
	if nInt > 0 {
		inv.SetCellValue(IN_SHEET, IN_DESC_INTERVIEWS, descInterview)
		inv.SetCellValue(IN_SHEET, IN_RATE_INTERVIEWS, rateInt)
		inv.SetCellValue(IN_SHEET, IN_UNITS_INTERVIEWS, nInt)
		inv.SetCellValue(IN_SHEET, IN_AMOUNT_INTERVIEWS, totalInterviews)
	}
	inv.SetCellValue(IN_SHEET, IN_AMOUNT, totalUsd)
	inv.SetCellValue(IN_SHEET, IN_TOTAL, totalOverall)

	db.Where("Month=?", month).Delete(&models.Invoice{})

	invoiceRecord := models.Invoice{
		Month:  month,
		Year:   year,
		Rate:   float64(rate),
		Hours:  float64(t),
		Amount: float64(totalUsd),
	}
	db.Create(&invoiceRecord)

	fmt.Printf("Saving invoice to %s\n", path+"/"+invFn)
	if err := inv.SaveAs(path + "/" + invFn); err != nil {
		log.Fatal(err)
	}
	fmt.Printf(
		"In month %d, %d hours were billed at %d USD/hour, total amount is %d USD\n",
		month,
		t,
		rate,
		totalUsd,
	)
	if nInt > 0 {
		fmt.Printf(
			"In month %d, %d interviews were billed at %d USD/interview, total amount is %d USD\n",
			month,
			nInt,
			rateInt,
			totalInterviews,
		)
		fmt.Printf("Overall amount is %d USD\n", totalOverall)
	}
	return totalUsd
}
