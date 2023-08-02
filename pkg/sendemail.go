package pkg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"github.com/timtarusov/timesheet_autofill/models"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func SendEmail(year int, month int, db *gorm.DB) {
	host := viper.GetString("Email.Sender.Host")
	port := viper.GetInt("Email.Sender.Port")
	user := viper.GetString("Email.Sender.Username")
	pwd := viper.GetString("Email.Sender.Password")
	from := viper.GetString("Email.Sender.Email")
	to := viper.GetString("Email.Recipient.Email")
	def_path := viper.GetString("Template.Default")
	timesheet_fn := viper.GetString("Template.TimesheetFilename")
	invoice_fn := viper.GetString("Template.InvoiceFilename")

	d := gomail.NewDialer(host, port, user, pwd)
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	month_s := time.Month(month).String()
	mon_map := map[int]string{
		1:  "JAN",
		2:  "FEB",
		3:  "MAR",
		4:  "APR",
		5:  "MAY",
		6:  "JUN",
		7:  "JUL",
		8:  "AUG",
		9:  "SEP",
		10: "OCT",
		11: "NOV",
		12: "DEC",
	}
	month_s_short := mon_map[month]
	year_s_short := strconv.Itoa(year)[2:]
	subject := fmt.Sprintf("Invoice for %s %d", month_s, year)
	m.SetHeader("Subject", subject)

	var inv = models.Invoice{}
	db.Where("Year=?", year).Where("Month=?", month).First(&inv)

	m.SetBody("text/html", CreateEmailBody(subject, year, month, inv.Hours, inv.Rate, inv.Amount))

	ts_path := fmt.Sprintf("%s/%s%s/%s", def_path, month_s_short, year_s_short, timesheet_fn)
	inv_path := fmt.Sprintf("%s/%s%s/%s", def_path, month_s_short, year_s_short, invoice_fn)

	m.Attach(ts_path)
	m.Attach(inv_path)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
