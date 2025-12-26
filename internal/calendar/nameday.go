package calendar

import (
	"context"
	"errors"
	"time"
)

type Nameday struct {
	ID   int
	Name string
	Date BirthdayDate
}

func NewNameday(
	id int,
	name string,
	date BirthdayDate,
) *Nameday {
	return &Nameday{
		ID:   id,
		Name: name,
		Date: date,
	}
}

var (
	ErrDuplicateNameday = errors.New("duplicate name day")
)

type NamedayStorage interface {
	// Save a Nameday
	//
	// Return an error of type duplicate if the same nameday or a similar with the same date already exists.
	Save(ctx context.Context, nd *Nameday) error

	// GetForDate returns the nameday for the specified month and day, if found.
	GetForDate(ctx context.Context, month time.Month, day uint8) (*Nameday, bool, error)

	// GetAll returns all the namedays.
	GetAll(ctx context.Context, limit uint, offset int) ([]*Nameday, error)

	// Delete a nameday given its ID and returns true if the target existed, false otherwise.
	// A non-nil error indiciates something went really wrong and the result might not be relevant.
	Delete(ctx context.Context, id int) (bool, error)
}
