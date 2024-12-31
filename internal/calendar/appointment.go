package calendar

import (
	"context"
	"errors"
	"time"
)

type Appointment struct {
	ID   int
	Name string

	// UTC
	Date time.Time
}

func NewAppointment(id int, name string, date time.Time) *Appointment {
	return &Appointment{
		ID:   id,
		Name: name,
		Date: date,
	}
}

var (
	ErrDuplicateAppointment = errors.New("duplicate appointment")
)

type AppointmentStorage interface {
	// Save an appointment
	//
	// If the same appointment already exists, an error of type
	// dupicate is returned
	Save(a *Appointment) error

	GetAllForDate(d uint8, m time.Month, y uint) ([]*Appointment, error)

	GetAll(limit uint8, offset int) ([]*Appointment, error)

	// Delete an appointment given its ID and returns true if the appointment
	// existed, false otherwise.
	// A non-nil error indiciates something went really wrong and the result
	// might not be relevant.
	Delete(ctx context.Context, id int) (bool, error)

	// DeleteAllBefore deletes all appointments before the provided date (time is ignored).
	// It returns the number of appointments deleted. If a non-nil error is returned,
	// the result must not be used.
	DeleteAllBefore(ctx context.Context, datetime time.Time) (int, error)
}
