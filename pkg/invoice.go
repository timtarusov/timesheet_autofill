package pkg

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/timtarusov/timesheet_autofill/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const (
	IN_SHEET  = "Simple Invoice"
	IN_DATE   = "E4"
	IN_NUM    = "E5"
	IN_DESC   = "A20"
	IN_UNITS  = "C20"
	IN_RATE   = "D20"
	IN_AMOUNT = "E20"
	IN_TOTAL  = "E37"
)

func Write_invoice(t int, year int, month int, path string, rate int, db *gorm.DB) int {
	templ_invoice := viper.GetString("Template.InvoicePath")
	inv_fn := viper.GetString("Template.InvoiceFilename")
	inv, err := excelize.OpenFile(templ_invoice)
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

	total_usd := t * rate

	desc_d_s := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	desc_d_e := desc_d_s.AddDate(0, 1, -1)
	desc := fmt.Sprintf("Consulting Servives - for the period of %s - %s", desc_d_s.Format("Jan 2, 2006"), desc_d_e.Format("Jan 2, 2006"))

	inv.SetCellValue(IN_SHEET, IN_DATE, d)
	inv.SetCellValue(IN_SHEET, IN_NUM, numInt+1)
	inv.SetCellInt(IN_SHEET, IN_UNITS, t)
	inv.SetCellValue(IN_SHEET, IN_DESC, desc)
	inv.SetCellValue(IN_SHEET, IN_RATE, rate)
	inv.SetCellValue(IN_SHEET, IN_AMOUNT, total_usd)
	inv.SetCellValue(IN_SHEET, IN_TOTAL, total_usd)

	db.Where("Month=?", month).Delete(&models.Invoice{})

	invoiceRecord := models.Invoice{
		Month:  month,
		Year:   year,
		Rate:   float64(rate),
		Hours:  float64(t),
		Amount: float64(total_usd),
	}
	db.Create(&invoiceRecord)

	fmt.Printf("Saving invoice to %s\n", path+"/"+inv_fn)
	if err := inv.SaveAs(path + "/" + inv_fn); err != nil {
		log.Fatal(err)
	}
	return total_usd

}
