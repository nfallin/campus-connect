package Model

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const enableForeignKeys = `PRAGMA foreign_keys=ON`

func Init(path string) (*sql.DB, error) {
	db, err := open(path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(enableForeignKeys); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	if _, err := db.Exec(usersTable); err != nil {
		return err
	}
	if _, err := db.Exec(eventsTable); err != nil {
		return err
	}
	if err := ensureEventInterestsColumn(db); err != nil {
		return err
	}
	if _, err := db.Exec(rsvpTable); err != nil {
		return err
	}
	return nil
}

func ensureEventInterestsColumn(db *sql.DB) error {
	rows, err := db.Query(`PRAGMA table_info(events)`)
	if err != nil {
		return err
	}
	hasInterests := false
	for rows.Next() {
		var columnID int
		var name, columnType string
		var notNull, primaryKey int
		var defaultValue sql.NullString
		if err := rows.Scan(&columnID, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			rows.Close()
			return err
		}
		if name == "interests" {
			hasInterests = true
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if hasInterests {
		return nil
	}
	_, err = db.Exec(`ALTER TABLE events ADD COLUMN interests TEXT NOT NULL DEFAULT '[]'`)
	return err
}
