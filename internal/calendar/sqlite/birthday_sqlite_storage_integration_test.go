//go:build integration

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/maximekuhn/metego/internal/calendar"
)

func TestSaveBd(t *testing.T) {
	// setup
	dbFileName := "test_db_save.sqlite3"
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

	sut, err := NewSQLiteBirthdayStorage(db)
	if err != nil {
		t.Errorf("failed to create sqlite birthday storage: %s", err)
		return
	}

	// actual test
	bd := calendar.NewBirthday(0, "Toto", time.February, 1)
	if err = sut.Save(bd); err != nil {
		t.Errorf("failed to save birthday: %s", err)
		return
	}
}

func TestSaveBdDuplicate(t *testing.T) {
	// setup
	dbFileName := "test_db_save.sqlite3"
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

	sut, err := NewSQLiteBirthdayStorage(db)
	if err != nil {
		t.Errorf("failed to create sqlite birthday storage: %s", err)
		return
	}

	// actual test
	bd := calendar.NewBirthday(0, "Toto", time.February, 1)
	if err = sut.Save(bd); err != nil {
		t.Errorf("failed to save birthday: %s", err)
		return
	}
	err = sut.Save(bd)
	if err == nil {
		t.Error("expected an error")
		return
	}

	expectedErrMsg := "duplicate birthday"
	actualErrMsg := err.Error()
	if expectedErrMsg != actualErrMsg {
		t.Errorf("want %s got %s", expectedErrMsg, actualErrMsg)
		return
	}
}

func TestGetAll(t *testing.T) {
	// setup
	dbFileName := "test_db_save.sqlite3"
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

	sut, err := NewSQLiteBirthdayStorage(db)
	if err != nil {
		t.Errorf("failed to create sqlite birthday storage: %s", err)
		return
	}

	// save some bd
	birthdays := []*calendar.Birthday{
		calendar.NewBirthday(0, "Toto", time.April, 2),
		calendar.NewBirthday(0, "Tata", time.May, 2),
		calendar.NewBirthday(0, "Doggo", time.April, 2),
	}
	for _, bd := range birthdays {
		err = sut.Save(bd)
		if err != nil {
			t.Errorf("failed to save bday: %s", err)
			return
		}
	}

	m := time.April
	d := 2
	allBds, err := sut.GetAllForDate(m, uint8(d))
	if err != nil {
		t.Errorf("failed to get all birthdays: %s", err)
		return
	}
	if len(allBds) != 2 {
		t.Errorf("want 2 birthdays, got %d", len(allBds))
		return
	}

	if !contains("Toto", allBds) {
		t.Error("expected Toto to be there")
		return
	}

	if !contains("Doggo", allBds) {
		t.Error("expected Doggo to be there")
		return
	}
}

func TestDeleteBday(t *testing.T) {
	// setup
	dbFileName := "test_db_save.sqlite3"
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

	sut, err := NewSQLiteBirthdayStorage(db)
	if err != nil {
		t.Errorf("failed to create sqlite birthday storage: %s", err)
		return
	}

	// save some bd
	birthdays := []*calendar.Birthday{
		calendar.NewBirthday(0, "Toto", time.April, 2),
		calendar.NewBirthday(0, "Tata", time.May, 2),
		calendar.NewBirthday(0, "Doggo", time.April, 2),
	}
	for _, bd := range birthdays {
		err = sut.Save(bd)
		if err != nil {
			t.Errorf("failed to save bday: %s", err)
			return
		}
	}

	// try to delete first birthday (id==1)
	found, err := sut.Delete(context.TODO(), 1)
	if err != nil {
		t.Fatalf("Delete(): expected ok got err %v", err)
	}
	if !found {
		t.Fatalf("Delete(): expected to find bday but not found")
	}

	// try to delete it again
	found, err = sut.Delete(context.TODO(), 1)
	if err != nil {
		t.Fatalf("Delete() again: expected ok got err %v", err)
	}
	if found {
		t.Fatalf("Delete() again: expected bday to have been deleted before but found it")
	}

	// try to delete a non-existing bday
	found, err = sut.Delete(context.TODO(), 999)
	if err != nil {
		t.Fatalf("Delete() non-existing: expected ok got err %v", err)
	}
	if found {
		t.Fatalf("Delete() non-existing: expected bday to not exist but found it")
	}
}

func contains(name string, bds []*calendar.Birthday) bool {
	for _, bd := range bds {
		if bd.Name == name {
			return true
		}
	}
	return false
}

func setupTmpDB(filePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, errors.New(fmt.Sprintf("could not ping db: %s", err.Error()))
	}

	if err = ApplyMigrations(db); err != nil {
		return nil, errors.New(fmt.Sprintf("could not apply migrations: %s", err.Error()))
	}

	return db, nil
}
