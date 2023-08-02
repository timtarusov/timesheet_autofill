package pkg

import (
	"fmt"
	"strconv"
	"time"

	"github.com/matcornic/hermes/v2"
	"github.com/spf13/viper"
)

func CreateEmailBody(subject string, year int, month int, hours float64, rate float64, amount float64) string {

	rec_name := viper.GetString("Email.Recipient.Name")
	sender_name := viper.GetString("Email.Sender.Name")
	sender_email := viper.GetString("Email.Sender.Email")
	h := hermes.Hermes{
		Product: hermes.Product{
			Name:        "Timofey",
			Link:        "tim.tarusov@yandex.ru",
			Copyright:   "Check out my github at https://github.com/timtarusov",
			TroubleText: "",
		},
		Theme: new(hermes.Flat),
	}

	month_s := time.Month(month).String()

	email := hermes.Email{
		Body: hermes.Body{
			Name:   rec_name,
			Intros: []string{fmt.Sprintf("Please find attached the timesheets and invoice for %s %d", month_s, year)},
			Table: hermes.Table{
				Data: [][]hermes.Entry{
					{
						{Key: "Name", Value: sender_name},
						{Key: "Year", Value: strconv.Itoa(year)},
						{Key: "Month", Value: month_s},
						{Key: "Hours", Value: fmt.Sprintf("%.0f", hours)},
						{Key: "Rate", Value: fmt.Sprintf("$ %.2f", rate)},
						{Key: "Amount", Value: fmt.Sprintf("$ %.2f", amount)},
					},
				},
				Columns: hermes.Columns{
					CustomWidth: map[string]string{
						"Name":   "30%",
						"Amount": "20%",
					},
					CustomAlignment: map[string]string{
						"Rate":   "right",
						"Amount": "right",
					},
				},
			},
			Actions: []hermes.Action{
				{
					Instructions: "Please confirm accepting this letter",
					Button: hermes.Button{
						Text: "Send letter of confirmation",
						Link: fmt.Sprintf("mailto:%s?subject=%s&body=Confirmed", sender_email, subject),
					},
				},
			},
			Signature: "Best regards",
		},
	}
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		panic(err) // Tip: Handle error with something else than a panic ;)
	}
	return emailBody

}
