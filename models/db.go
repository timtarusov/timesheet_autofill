package models

import "gorm.io/gorm"

type Timesheet struct {
	gorm.Model
	Day int
	Month int
	Year int
	Value int
}

type Invoice struct {
	gorm.Model
	Month int
	Year int
	Rate float64
	Hours float64
	Amount float64
}