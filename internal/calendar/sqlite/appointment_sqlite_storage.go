package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/internal/calendar"
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
func (s *SQLiteAppointmentStorage) Save(ctx context.Context, a *calendar.Appointment) error {
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

func (s *SQLiteAppointmentStorage) GetAllForDate(ctx context.Context, d uint8, m time.Month, y uint) ([]*calendar.Appointment, error) {
	rows, err := s.db.Query(
		"SELECT id, name, date FROM appointments WHERE date >= ?",
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

func (s *SQLiteAppointmentStorage) GetAll(ctx context.Context, limit uint8, offset int) ([]*calendar.Appointment, error) {
	rows, err := s.db.Query(
		"SELECT id, name, date FROM appointments LIMIT ? OFFSET ?",
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
		var id int
		var name string
		var dateStr string
		if err := rows.Scan(&id, &name, &dateStr); err != nil {
			return nil, err
		}

		if id < 1 {
			return nil, fmt.Errorf("expected ID > 1 but got %d", id)
		}

		// keep only date, ignore time
		date, err := time.Parse("2006-01-02 15:04:05-07:00", dateStr)
		if err != nil {
			return nil, err
		}

		apt := calendar.NewAppointment(id, name, date)
		apts = append(apts, apt)
	}
	return apts, nil
}

func (s *SQLiteAppointmentStorage) Delete(ctx context.Context, id int) (bool, error) {
	query := "DELETE FROM appointments WHERE id = ?"
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, err
	}
	affectedCount, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affectedCount == 1, nil
}

func (s *SQLiteAppointmentStorage) DeleteAllBefore(ctx context.Context, date time.Time) (int, error) {
	query := "DELETE FROM appointments WHERE date < ?"
	res, err := s.db.ExecContext(ctx, query, fmt.Sprintf("%04d-%02d-%02d", date.Year(), date.Month(), date.Day()))
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}
