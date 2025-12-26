package sqlite

import "database/sql"

func ApplyMigrations(db *sql.DB) error {
	err := createBirhtdaysTable(db)
	if err != nil {
		return err
	}

	err = createAppointmentsTable(db)
	if err != nil {
		return err
	}

	err = createNamedaysTable(db)
	if err != nil {
		return err
	}

	return nil
}

func createBirhtdaysTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS birthdays (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        date TEXT,
        UNIQUE (name, date)
    )
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func createAppointmentsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS appointments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        date TEXT,
        UNIQUE (name, date)
    )
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func createNamedaysTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS namedays (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT INTEGER NOT NULL UNIQUE,
        month INTEGER NOT NULL,
		day INTEGER NOT NULL,
        UNIQUE (month, day)
    )
    `
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
