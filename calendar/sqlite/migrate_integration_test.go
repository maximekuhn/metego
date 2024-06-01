//go:build integration

package sqlite

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	dbFileName := "test_db_apply_migrations.sqlite3"
	f, err := os.CreateTemp(".", dbFileName)
	if err != nil {
		t.Errorf("failed to create tmp db file: %s", err)
		return
	}
    defer os.Remove(f.Name())

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		t.Errorf("failed to create db driver: %s", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		t.Errorf("failed to ping db: %s", err)
		return
	}

	if err = ApplyMigrations(db); err != nil {
		t.Errorf("failed to apply migrations: %s", err)
		return
	}
}
