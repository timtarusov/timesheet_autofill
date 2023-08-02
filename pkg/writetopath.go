package pkg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func Create_path(year int, month int) (string, error) {
	def_path := viper.GetString("Template.Default")
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

	path := def_path + "/" + mon_map[month] + strconv.Itoa(year)[2:]
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	fmt.Printf("New Path: %s\n", path)
	return path, nil
}
