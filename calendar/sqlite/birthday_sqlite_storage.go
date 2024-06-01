package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/calendar"
)

type SQLiteBdayStorage struct {
	db *sql.DB
}

func NewSQliteBdayStorage(db *sql.DB) (*SQLiteBdayStorage, error) {
	return &SQLiteBdayStorage{
		db: db,
	}, nil
}

// save the birthday
// if the same birthday already exists, an error is returned
func (s *SQLiteBdayStorage) Save(b *calendar.Birthday) error {
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
func (s *SQLiteBdayStorage) GetAllForDate(month time.Month, day uint8) ([]*calendar.Birthday, error) {
	rows, err := s.db.Query(
		"SELECT name, date FROM birthdays WHERE date = ?",
		fmt.Sprintf("%d/%d", month, day),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bdays := make([]*calendar.Birthday, 0)
	for rows.Next() {
		var name string
		var dateStr string
		if err := rows.Scan(&name, &dateStr); err != nil {
			return nil, err
		}

		parts := strings.Split(dateStr, "/")
		if len(parts) != 2 {
			return nil, errors.New(fmt.Sprintf("date is corrupted for %s", name))
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

		bday := calendar.NewBirthday(name, time.Month(m), uint8(d))
		bdays = append(bdays, bday)
	}

	return bdays, nil
}
