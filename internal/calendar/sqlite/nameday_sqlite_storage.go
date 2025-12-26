package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/internal/calendar"
)

type SQLiteNamedayStorage struct {
	db *sql.DB
}

func NewSQLiteNamedayStorage(db *sql.DB) *SQLiteNamedayStorage {
	return &SQLiteNamedayStorage{
		db: db,
	}
}

func (s *SQLiteNamedayStorage) Save(ctx context.Context, nd *calendar.Nameday) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO namedays
		(name, month, day) VALUES
		(?, ?, ?)
		`, nd.Name, nd.Date.Month, nd.Date.Day)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return calendar.ErrDuplicateNameday
			}
		}
		return err
	}
	return nil
}

func (s *SQLiteNamedayStorage) GetForDate(ctx context.Context, month time.Month, day uint8) (*calendar.Nameday, bool, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name
		FROM namedays
		WHERE month = ?
		AND day = ?
		`, month, day)

	var id int
	var name string

	if err := row.Scan(&id, &name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return calendar.NewNameday(id, name, calendar.BirthdayDate{
		Month: month,
		Day:   day,
	}), true, nil
}

func (s *SQLiteNamedayStorage) GetAll(ctx context.Context, limit uint, offset int) ([]*calendar.Nameday, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, month, day
		FROM namedays
		ORDER BY month, day
		LIMIT ? OFFSET ?
		`, limit, offset)

	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	namedays := make([]*calendar.Nameday, 0)
	for rows.Next() {
		var id int
		var name string
		var month int
		var day uint8

		if err := rows.Scan(&id, &name, &month, &day); err != nil {
			return namedays, err
		}
		namedays = append(namedays, calendar.NewNameday(id, name, calendar.BirthdayDate{
			Month: time.Month(month),
			Day:   day,
		}))
	}
	return namedays, nil
}

func (s *SQLiteNamedayStorage) Delete(ctx context.Context, id int) (bool, error) {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM namedays WHERE id = ?
		`, id)
	if err != nil {
		return false, err
	}
	affectedCount, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affectedCount == 1, nil
}
