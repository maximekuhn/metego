package calendar

import (
	"errors"
	"time"
)

type Appointment struct {
	Name string

	// UTC
	Date time.Time
}

func NewAppointment(name string, date time.Time) *Appointment {
	return &Appointment{
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

	GetAllForDate(d uint8, m time.Month, y uint8) ([]*Appointment, error)

	GetAll(limit uint8, offset int) ([]*Appointment, error)
}
