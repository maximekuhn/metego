package sqlite

import "database/sql"

func ApplyMigrations(db *sql.DB) error {
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
