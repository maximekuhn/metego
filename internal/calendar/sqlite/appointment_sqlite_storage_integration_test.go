//go:build integration

package sqlite

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/maximekuhn/metego/internal/calendar"
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
	apt := calendar.NewAppointment(0, "car", time.Now())
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
	apt := calendar.NewAppointment(0, "car", time.Now())
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

func TestGetAllApts(t *testing.T) {
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
	apts, err := sut.GetAll(10, 0)
	if err != nil {
		t.Fatalf("GetAll(): expected ok got %v", err)
	}
	if len(apts) != 10 {
		t.Fatalf("GetAll(): expected to retrieve all 10 appointments, got %d", len(apts))
	}

	for _, apt := range apts {
		if apt.ID < 1 {
			t.Errorf("GetAll(): ID should never be < 1 but got %d (apt name: %s)", apt.ID, apt.Name)
		}
	}
}

func TestDeleteApt(t *testing.T) {
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

	// scenarios
	testcases := []struct {
		title string
		// fixtures generate 10 appointments, with ID starting at 1
		aptID                int
		aptShouldHaveExisted bool
	}{
		{
			title:                "non existing apt",
			aptID:                25,
			aptShouldHaveExisted: false,
		},
		{
			title:                "existing apt",
			aptID:                2,
			aptShouldHaveExisted: true,
		},
	}

	for _, test := range testcases {
		t.Run(test.title, func(t *testing.T) {
			found, err := sut.Delete(context.TODO(), test.aptID)
			if err != nil {
				t.Fatalf("Delete(%d): expected ok got err %v", test.aptID, err)
			}
			if test.aptShouldHaveExisted && !found {
				t.Fatalf("Delete(%d): expected appointment to exist", test.aptID)
			}

			if !test.aptShouldHaveExisted && found {
				t.Fatalf("Delete(%d): expected no appointment to exist with this ID", test.aptID)
			}

			apts, err := sut.GetAll(10, 0)
			for _, apt := range apts {
				if apt.ID == test.aptID {
					t.Fatalf("Delete(%d): expected apt to be deleted but it's still here", test.aptID)
				}
			}
		})
	}
}

func TestDeleteAllAptsBefore(t *testing.T) {
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

	// add 2 apts
	date, err := time.Parse("2006-01-02", "2023-07-14")
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}
	if err := sut.Save(calendar.NewAppointment(0, "Football match", date)); err != nil {
		t.Fatalf("Failed to save apt: %v", err)
	}
	if err := sut.Save(calendar.NewAppointment(0, "MOTOGP race", date)); err != nil {
		t.Fatalf("Failed to save apt: %v", err)
	}
	if err := sut.Save(calendar.NewAppointment(0, "Job interview", date.Add(2*24*time.Hour))); err != nil {
		t.Fatalf("Failed to save apt: %v", err)
	}

	limitDate := date.Add(24 * time.Hour)
	deleted, err := sut.DeleteAllBefore(context.TODO(), limitDate)
	if err != nil {
		t.Fatalf("DeleteAllBefore(): expected ok got err %v", err)
	}
	if deleted != 2 {
		t.Fatalf("DeleteAllBefore(): expected to delete 2 apts but deleted %d", deleted)
	}
}

func fixturesAppoitments(sut *SQLiteAppointmentStorage) error {
	for i := 1; i <= 10; i++ {
		name := fmt.Sprintf("Appointment %d", i)
		date, _ := time.Parse("2006-01-02", fmt.Sprintf("2023-07-%02d", i))
		err := sut.Save(calendar.NewAppointment(0, name, date))
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
