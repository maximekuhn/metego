package calendar

import (
	"errors"
	"time"
)

var (
	ErrDuplicateBirthday = errors.New("duplicate birthday")
)

type BirthdayStorage interface {
	// save the birthday
	// if the same birthday already exists, an error is returned
	Save(b *Birthday) error

	// get all birthdays given the current date
	GetAllForDate(month time.Month, day uint8) ([]*Birthday, error)

    // get all birthdays
    GetAll(limit, offset int) ([]*Birthday, error)
}
