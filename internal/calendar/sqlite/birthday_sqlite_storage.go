package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/internal/calendar"
)

type SQLiteBirthdayStorage struct {
	db *sql.DB
}

func NewSQLiteBirthdayStorage(db *sql.DB) (*SQLiteBirthdayStorage, error) {
	return &SQLiteBirthdayStorage{
		db: db,
	}, nil
}

// save the birthday
// if the same birthday already exists, an error is returned
func (s *SQLiteBirthdayStorage) Save(ctx context.Context, b *calendar.Birthday) error {
	_, err := s.db.Exec(
		"INSERT INTO birthdays (name, date) VALUES (?, ?)",
		b.Name,
		fmt.Sprintf("%d/%d", b.Date.Month, b.Date.Day),
	)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return calendar.ErrDuplicateBirthday
			}
		}
		return err
	}

	return nil
}

// get all birthdays given the current date
func (s *SQLiteBirthdayStorage) GetAllForDate(ctx context.Context, month time.Month, day uint8) ([]*calendar.Birthday, error) {
	rows, err := s.db.Query(
		"SELECT id, name, date FROM birthdays WHERE date = ?",
		fmt.Sprintf("%d/%d", month, day),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bdays, err := convertRowsBdays(rows)
	if err != nil {
		return nil, err
	}

	return bdays, nil
}

func (s *SQLiteBirthdayStorage) GetAll(ctx context.Context, limit uint, offset int) ([]*calendar.Birthday, error) {
	rows, err := s.db.Query(
		"SELECT id, name, date FROM birthdays LIMIT ? OFFSET ?",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bdays, err := convertRowsBdays(rows)
	if err != nil {
		return nil, err
	}

	return bdays, nil
}

func (s *SQLiteBirthdayStorage) Delete(ctx context.Context, id int) (bool, error) {
	query := "DELETE FROM birthdays WHERE id = ?"
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

func convertRowsBdays(rows *sql.Rows) ([]*calendar.Birthday, error) {
	bdays := make([]*calendar.Birthday, 0)
	for rows.Next() {
		var id int
		var name string
		var date string
		if err := rows.Scan(&id, &name, &date); err != nil {
			return nil, err
		}

		if id < 1 {
			return nil, fmt.Errorf("expected ID to be >=1 but got %d", id)
		}

		parts := strings.Split(date, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("date is corrupted for %s", name)
		}

		month := parts[0]
		day := parts[1]

		m, err := strconv.ParseInt(month, 10, 64)
		if err != nil {
			return nil, err
		}

		d, err := strconv.ParseInt(day, 10, 8)
		if err != nil {
			return nil, err
		}

		bday := calendar.NewBirthday(id, name, time.Month(m), uint8(d))
		bdays = append(bdays, bday)
	}
	return bdays, nil
}
