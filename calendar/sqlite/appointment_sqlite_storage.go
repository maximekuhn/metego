package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/calendar"
)

type SQLiteAppointmentStorage struct {
	db *sql.DB
}

func NewSQLiteAppointmentStorage(db *sql.DB) *SQLiteAppointmentStorage {
	return &SQLiteAppointmentStorage{
		db: db,
	}
}

// Save an appointment
// If the same appointment already exists, an error of type
// dupicate is returned
func (s *SQLiteAppointmentStorage) Save(a *calendar.Appointment) error {
	_, err := s.db.Exec(
		"INSERT INTO appointments (name, date) VALUES (?,?)",
		a.Name,
		a.Date,
	)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return calendar.ErrDuplicateAppointment
			}
		}
		return err
	}

	return nil
}

func (s *SQLiteAppointmentStorage) GetAllForDate(d uint8, m time.Month, y uint) ([]*calendar.Appointment, error) {
	rows, err := s.db.Query(
		"SELECT name, date FROM appointments WHERE date >= ?",
		fmt.Sprintf("%04d-%02d-%02d", y, m, d),
	)
	if err != nil {
		return nil, err
	}

	apts, err := convertRowsApts(rows)
	if err != nil {
		return nil, err
	}

	// TODO: fix query so we don't have to do it manually, as it
	// is very sub-optimal
	appointments := make([]*calendar.Appointment, 0)
	for _, apt := range apts {
		if uint(apt.Date.Year()) != y {
			continue
		}

		if apt.Date.Month() != m {
			continue
		}

		if uint8(apt.Date.Day()) != d {
			continue
		}

		appointments = append(appointments, apt)
	}

	return appointments, nil
}

func (s *SQLiteAppointmentStorage) GetAll(limit uint8, offset int) ([]*calendar.Appointment, error) {
	rows, err := s.db.Query(
		"SELECT name, date FROM appointments LIMIT ? OFFSET ?",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apts, err := convertRowsApts(rows)
	if err != nil {
		return nil, err
	}

	return apts, nil
}

func convertRowsApts(rows *sql.Rows) ([]*calendar.Appointment, error) {
	apts := make([]*calendar.Appointment, 0)
	for rows.Next() {
		var name string
		var dateStr string
		if err := rows.Scan(&name, &dateStr); err != nil {
			return nil, err
		}

		// keep only date, ignore time
		date, err := time.Parse("2006-01-02 15:04:05-07:00", dateStr)
		if err != nil {
			return nil, err
		}

		apt := calendar.NewAppointment(name, date)
		apts = append(apts, apt)
	}
	return apts, nil
}
