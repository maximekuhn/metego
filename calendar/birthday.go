package calendar

import (
	"time"
)

type Birthday struct {
	Name string
	Date BirthdayDate
}

type BirthdayDate struct {
	Month time.Month
	Day   uint8
}

func NewBirthday(name string, m time.Month, d uint8) *Birthday {
	return &Birthday{
		Name: name,
		Date: BirthdayDate{
			Month: m,
			Day:   d,
		},
	}
}
