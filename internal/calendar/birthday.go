package calendar

import (
	"context"
	"errors"
	"time"
)

type Birthday struct {
	ID   int
	Name string
	Date BirthdayDate
}

type BirthdayDate struct {
	Month time.Month
	Day   uint8
}

func NewBirthday(id int, name string, m time.Month, d uint8) *Birthday {
	return &Birthday{
		ID:   id,
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

	// Delete a birthday given its ID and returns true if the target existed, false otherwise.
	// A non-nil error indiciates something went really wrong and the result might not be relevant.
	Delete(ctx context.Context, id int) (bool, error)
}
