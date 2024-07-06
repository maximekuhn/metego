//go:build integration

package sqlite

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/maximekuhn/metego/calendar"
)

func TestSaveApt(t *testing.T) {
	// setup
	dbFileName := "test_db_save_apt.sqlite3"
	f, err := os.CreateTemp(".", dbFileName)
	if err != nil {
		t.Fatalf("failed to create tmp db file: %s", err)
		return
	}
	defer os.Remove(f.Name())

	db, err := setupTmpDB(f.Name())
	if err != nil {
		t.Errorf("failed to create db driver: %s", err)
		return
	}
	defer db.Close()

	sut := NewSQLiteAppointmentStorage(db)

	// actual test
	apt := calendar.NewAppointment("car", time.Now())
	if err = sut.Save(apt); err != nil {
		t.Errorf("failed to save appointment: %s", err)
		return
	}

}

func TestSaveAptDuplicate(t *testing.T) {
	// setup
	dbFileName := "test_db_save_apt.sqlite3"
	f, err := os.CreateTemp(".", dbFileName)
	if err != nil {
		t.Fatalf("failed to create tmp db file: %s", err)
		return
	}
	defer os.Remove(f.Name())

	db, err := setupTmpDB(f.Name())
	if err != nil {
		t.Errorf("failed to create db driver: %s", err)
		return
	}
	defer db.Close()

	sut := NewSQLiteAppointmentStorage(db)

	// actual test
	apt := calendar.NewAppointment("car", time.Now())
	_ = sut.Save(apt)
	err = sut.Save(apt)
	if err == nil {
		t.Error("should have gotten an error")
	}
	if !errors.Is(calendar.ErrDuplicateAppointment, err) {
		t.Errorf("expected ErrDuplicateAppointment got '%s'", err)
	}
}

func TestGetAllFor(t *testing.T) {
	// setup
	dbFileName := "test_db_save_apt.sqlite3"
	f, err := os.CreateTemp(".", dbFileName)
	if err != nil {
		t.Fatalf("failed to create tmp db file: %s", err)
		return
	}
	defer os.Remove(f.Name())

	db, err := setupTmpDB(f.Name())
	if err != nil {
		t.Errorf("failed to create db driver: %s", err)
		return
	}
	defer db.Close()

	sut := NewSQLiteAppointmentStorage(db)

	// setup appoitments
	err = fixturesAppoitments(sut)
	if err != nil {
		t.Errorf("fixtures failed: %s", err)
	}

	// try to fetch
	date, _ := time.Parse("2006-01-02", "2023-07-04")
	apts, err := sut.GetAllForDate(uint8(date.Day()), date.Month(), uint(date.Year()))
	if err != nil {
		t.Fatalf("failed to get appointments: %s", err)
	}
	if !containsApt("Appointment 4", apts) {
		t.Error("Appointment 4 is not here")
	}
	if len(apts) != 1 {
		t.Errorf("got %d appointments, wanted 1", len(apts))
	}
}

func fixturesAppoitments(sut *SQLiteAppointmentStorage) error {
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("Appointment %d", i)
		date, _ := time.Parse("2006-01-02", fmt.Sprintf("2023-07-%02d", i))
		err := sut.Save(calendar.NewAppointment(name, date))
		if err != nil {
			return err
		}
	}
	return nil
}

func containsApt(name string, apts []*calendar.Appointment) bool {
	for _, apt := range apts {
		if apt.Name == name {
			return true
		}
	}
	return false
}
