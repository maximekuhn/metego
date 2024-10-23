package calendar

import (
	"errors"
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

var (
	ErrDuplicateBirthday = errors.New("duplicate birthday")
)

type BirhtdayStorage interface {
	// Save a birthday
	//
	// Return an error of type duplicate if the same birthday already exists
	Save(b *Birthday) error

	// GetAllForDate returns all the birthdays for the given date
	GetAllForDate(m time.Month, day uint8) ([]*Birthday, error)

	// GetAll returns all the birthdays
	GetAll(limit uint, offset int) ([]*Birthday, error)
}
